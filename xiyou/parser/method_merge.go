package parser

import (
	"bytes"
	"fmt"
	"strconv"
)

type MergeToObjectArray struct {
	OutputKey  string `mapstructure:"output_key"`
	OriginKeys []struct {
		Keys         []*BeMergedKey `mapstructure:"keys"`
		Position     int            `mapstructure:"position"`
		PositionName string         `mapstructure:"position_name"`
	} `mapstructure:"origin_keys"`
}

type BeMergedKey struct {
	Type      string `mapstructure:"type"`
	Name      string `mapstructure:"name"`
	OutputKey string `mapstructure:"output_key"`
}

// merge 合并多个key到一个对象组 template标头拆分模板
func merge(splitBytes [][]byte, maps []*MergeToObjectArray, template *setting) error {
	for _, cfg := range maps {
		bup := bytes.NewBuffer([]byte{})
		bup.WriteString("[")
		for i, ori := range cfg.OriginKeys {
			tmpBuf := bytes.NewBuffer([]byte{})
			tmpBuf.WriteString("{")
			tmpBuf.WriteString(fmt.Sprintf("\"%s\": %d", ori.PositionName, ori.Position))
			tmpStr := ""
			for _, key := range ori.Keys {
				continueMap[key.Name] = struct{}{}
				idx, ok := template.idxs[key.Name]
				if !ok {
					return fmt.Errorf("[fn merge] key[%s] not found", key.Name)
				}
				if idx > len(splitBytes)-1 {
					zz := ""
					for _, splitByte := range splitBytes {
						zz += string(splitByte)
					}
					return fmt.Errorf("[fn merge]request body[%s] and template[%s] not match", zz, template.Template)
				}
				content := splitBytes[idx]
				switch key.Type {
				case "int":
					var num int
					if string(content) == NULL {
						num = 0
					} else {
						tmp, err := strconv.Atoi(string(content))
						if err != nil {
							return fmt.Errorf("[fn merge] file[%s] [key:%s] [data:(%s)] can not convert to int", template.Filename, key.Name, string(content))
						}
						num = tmp
					}

					tmpStr = fmt.Sprintf(intFormat, key.OutputKey, num)
				case "float":
					var num float64
					if string(content) == NULL {
						num = 0
					} else {
						tmp, err := strconv.ParseFloat(string(content), 64)
						if err != nil {
							return fmt.Errorf("[fn merge][file:%s] [key:%s] [data:(%s)] can not convert to float64", template.Filename, key.Name, string(content))
						}
						num = tmp
					}
					tmpStr = fmt.Sprintf(floatFormat, key.OutputKey, num)
				default:
					tmpStr = fmt.Sprintf(stringFormat, key.OutputKey, string(content))
				}
				tmpBuf.WriteString(",")
				tmpBuf.WriteString(tmpStr)
			}
			tmpBuf.WriteString("}")
			if i > 0 {
				bup.WriteString(",")
			}
			bup.Write(tmpBuf.Bytes())
		}
		bup.WriteString("]")
		mergeProperties[cfg.OutputKey] = bup.String()
	}
	return nil
}
