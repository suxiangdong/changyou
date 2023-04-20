package parser2

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

var propertiesMap = map[string]string{}

// Parse 解析数据
func Parse(bs []byte) ([]byte, error) {
	buf := bufpool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufpool.Put(buf)
	}()
	//b := bytes.Replace(bs, []byte("\""), []byte("\\\""), -1)
	// 转义双引号和转义符
	rep := strings.NewReplacer("\"", "\\\"", "\\", "\\\\")
	st := rep.Replace(string(bs))
	b := []byte(st)

	dds := bytes.Split(b, []byte(separator))
	ia := 0
	for _, dd := range dds {
		if ia > 0 {
			buf.Write([]byte(separator))
		}
		c, err := parseLine(dd)
		if err != nil {
			return nil, err
		}
		buf.Write(c)
		ia++
	}
	return buf.Bytes(), nil
}

func parseLine(b []byte) ([]byte, error) {
	propertiesMap = map[string]string{}
	pb := bufpool.Get().(*bytes.Buffer)
	defer func() {
		pb.Reset()
		bufpool.Put(pb)
	}()
	// remove \n \r\n
	if b[len(b)-1] == '\n' {
		drop := 1
		if len(b) > 1 && b[len(b)-2] == '\r' {
			drop = 2
		}
		b = b[:len(b)-drop]
	}
	d := bytes.Split(b, []byte("|ta|"))
	if len(d) != 2 {
		return nil, fmt.Errorf("request data is invalid, [%s]", string(b))
	}
	// 整条数据为null，原样返回去，在logbus内报错
	if string(d[1]) == NULL {
		return b, nil
	}
	template, ok := cfg.Templates[string(d[0])]
	if !ok {
		return nil, fmt.Errorf("template [%s] not found", string(d[0]))
	}
	idx := 0
	for _, cfg := range template.Configs {
		if idx > 0 {
			pb.WriteByte(0x01)
		}
		idx++
		ix := 0
		pb.WriteByte('{')
		splitBytes := bytes.Split(d[1], []byte{0x01})
		existIdxs := map[int]struct{}{}
		for origin, taField := range cfg.Mapping {
			if cfg.Filter != nil {
				_, ok := cfg.Filter[origin]
				_, ok2 := standardTaFields[taField]
				if !ok && !ok2 {
					continue
				}
			}
			idx, ok := template.idxs[origin]
			if !ok {
				continue
			}
			zz := ""
			for _, bb := range splitBytes {
				zz += string(bb)
			}
			if idx > len(splitBytes)-1 {
				return nil, fmt.Errorf("[%d|%d|%s]][%s]request body[%s] and template[%s] not match", idx, len(splitBytes)-1, zz, template.Filename, string(b), template.Template)
			}
			if string(splitBytes[idx]) == NULL {
				continue
			}
			// 如果需要走特殊逻辑
			var str string
			var field *field
			if _, ok := cfg.Fields[origin]; ok {
				field = cfg.Fields[origin]
				if field.Idx != 9999 {
					idx = field.Idx
				}
				if field.Params != nil {
					var err error
					str, err = parseSpecialField(field, taField, splitBytes[idx], template.Filename)
					if err != nil {
						return nil, err
					}
				}
			}
			// 未经过特殊处理
			if str == "" {
				if field != nil {
					switch field.Type {
					case "float":
						num, err := strconv.ParseFloat(string(splitBytes[idx]), 64)
						if err != nil {
							return nil, fmt.Errorf("[%s]data[%s] can not convert to float64", template.Filename, string(splitBytes[idx]))
						}
						str = fmt.Sprintf(floatFormat, taField, num)
					case "int":
						num, err := strconv.Atoi(string(splitBytes[idx]))
						if err != nil {
							return nil, fmt.Errorf("[%s]data[%s] can not convert to int", template.Filename, string(splitBytes[idx]))
						}
						str = fmt.Sprintf(intFormat, taField, num)
					default:
					}
				}
				if str == "" {
					str = fmt.Sprintf(stringFormat, taField, string(splitBytes[idx]))
				}
			}
			if _, ok := standardTaFields[taField]; !ok {
				propertiesMap[origin] = str
				continue
			}
			if ix > 0 {
				pb.WriteByte(',')
			}
			existIdxs[idx] = struct{}{}
			pb.WriteString(str)
			ix++
		}
		typ := defaultType
		if cfg.Type != "" {
			typ = cfg.Type
		}
		defaultStr := fmt.Sprintf(",\"#type\":\"%s\"", typ)
		pb.WriteString(defaultStr)
		properties, err := parseProperties(splitBytes, existIdxs, template, cfg)
		if err != nil {
			return nil, err
		}
		if string(properties) != "" {
			pb.WriteString(",\"properties\":")
			pb.Write(properties)
		}
		pb.WriteByte('}')
	}

	return pb.Bytes(), nil
}

func parseSpecialField(f *field, taField string, originContext []byte, fname string) (string, error) {
	fn, ok := fieldHandleFns[f.Params.Fn]
	if !ok {
		return "", fmt.Errorf("[%s]handle fn[%s] not exist", fname, f.Params.Fn)
	}
	tempStr, err := fn(originContext, f.Params, false)
	if err != nil {
		return "", fmt.Errorf("[%s]%v", fname, err)
	}
	return fmt.Sprintf(format, taField, tempStr), nil
}

func parseProperties(b [][]byte, exists map[int]struct{}, template *setting, cfg *propertiesSetting) ([]byte, error) {
	pbf := bufpool.Get().(*bytes.Buffer)
	defer func() {
		pbf.Reset()
		bufpool.Put(pbf)
	}()
	pbf.WriteByte('{')
	ix := 0

	for i, v := range b {
		if string(v) == NULL {
			continue
		}
		oriKey := template.templates[i]
		if cfg.Filter != nil {
			_, ok := cfg.Filter[oriKey]
			if !ok {
				continue
			}
		}
		repeat := ""
		if _, ok := exists[i]; ok {
			if cfg.Repeat == nil {
				continue
			}
			temp, ok := cfg.Repeat[oriKey]
			if !ok {
				continue
			}
			repeat = temp
		}
		if ix > 0 {
			pbf.WriteByte(',')
		}
		var str string
		if i >= len(template.templates) {
			return []byte{}, fmt.Errorf("[%s] title and content num is not match", template.Filename)
		}
		if t, ok := propertiesMap[oriKey]; ok {
			str = t
		} else {
			k := template.templates[i]
			if repeat != "" {
				k = repeat
			}
			str = fmt.Sprintf("\"%s\":\"%s\"", k, string(v))
		}
		pbf.WriteString(str)
		ix++
	}
	pbf.WriteByte('}')
	return pbf.Bytes(), nil
}
