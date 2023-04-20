package parser2

import (
	"bytes"
	"fmt"
	"strconv"
)

const discard = 1
const typeInline = 1

var sep = map[string]int8{
	",": 0,
	"_": 0,
}
var objSep = map[string]int8{
	"|": 0,
}

type StringToOtherParams struct {
	SpecialOp       int `json:"specialOp"`
	Special         map[string]map[string]int
	Keys            []*Key          `json:"keys"`
	Fn              string          `json:"fn"`
	Sep             map[string]int8 `json:"sep"`    // 数组分隔符
	Objsep          map[string]int8 `json:"objsep"` // 对象分隔符
	Symbols         string          `json:"symbols"`
	Prefix          string          `json:"prefix"`
	AddSerialNumber bool            `json:"addSerialNumber"`
}

type Key struct {
	Name       string               `json:"name"`
	Type       string               `json:"type"`
	Handletype int                  `json:"handletype"`
	Params     *StringToOtherParams `json:"params"`
}

var fieldHandleFns = map[string]func(b []byte, params *StringToOtherParams, inline bool) (string, error){
	"string2objectarray": stringToObjectArray,
}

func bufWriteByte(buf *bytes.Buffer, b byte, in bool) {
	if !in {
		buf.WriteByte(b)
	}
}

func getSpecialIdx(m map[string]int, k string) (int, bool) {
	idx, ok := m[k]
	if !ok {
		idx, ok = m["all"]
	}
	return idx, ok
}

func stringToObjectArray(b []byte, params *StringToOtherParams, inline bool) (string, error) {
	// 没有参数就原样返回
	if params == nil {
		return string(b), nil
	}
	s, os := sep, objSep
	if params.Sep != nil {
		s = params.Sep
	}
	if params.Objsep != nil {
		os = params.Objsep
	}
	buf := bufpool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufpool.Put(buf)
	}()
	bufWriteByte(buf, '[', inline)
	ikx := 0
	for i, v := range bytes.FieldsFunc(b, func(r rune) bool {
		_, ok := s[string([]rune{r})]
		return ok
	}) {
		// NULL 跳过
		if string(v) == NULL {
			continue
		}
		op := bytes.FieldsFunc(v, func(r rune) bool {
			_, ok := os[string([]rune{r})]
			return ok
		})
		if len(op) <= 0 {
			continue
		}
		if len(op) > len(params.Keys) {
			return "", fmt.Errorf("string[%s] to objectArray[%v] length not match", string(v), params.Keys)
		}
		ikx++
		opless := len(op) < len(params.Keys)
		if opless && params.SpecialOp == discard {
			continue
		}
		if opless && params.Special == nil {
			return "", fmt.Errorf("special is not setting")
		}
		if i > 0 {
			buf.WriteByte(',')
		}
		bufWriteByte(buf, '{', inline)
		for ik, key := range params.Keys {
			if key.Handletype == typeInline {
				str, err := stringToObjectArray(op[ik], key.Params, true)
				if err != nil {
					return "", err
				}
				if str != "" {
					buf.WriteByte(',')
				}
				buf.WriteString(str)
				continue
			}
			if ik > 0 {
				buf.WriteByte(',')
			}
			var str string
			keyName := key.Name
			if inline {
				if params.AddSerialNumber {
					keyName = fmt.Sprintf("%s%s%d%s%s", params.Prefix, params.Symbols, ikx, params.Symbols, key.Name)
					if key.Name == "" {
						keyName = fmt.Sprintf("%s%s%d", params.Prefix, params.Symbols, ikx)
					}
				} else {
					keyName = fmt.Sprintf("%s%s%s", params.Prefix, params.Symbols, key.Name)
				}
			}
			switch key.Type {
			case "int":
				if opless {
					tem, ok := params.Special[string(v)]
					if !ok {
						return "", fmt.Errorf("special is not matching")
					}
					if idx, ok := getSpecialIdx(tem, key.Name); ok {
						num, err := strconv.Atoi(string(op[idx]))
						if err != nil {
							return "", fmt.Errorf("data[%s] can not convert to int", string(op[idx]))
						}
						str = fmt.Sprintf(intFormat, keyName, num)
					} else {
						str = fmt.Sprintf(intFormat, keyName, 0)
					}
				} else {
					num, err := strconv.Atoi(string(op[ik]))
					if err != nil {
						return "", fmt.Errorf("data[%s] can not convert to int", string(op[ik]))
					}
					str = fmt.Sprintf(intFormat, keyName, num)
				}
			case "float":
				if inline {
					keyName = fmt.Sprintf("%s%s%d%s%s", params.Prefix, params.Symbols, ik+1, params.Symbols, key.Name)
				}
				if opless {
					tem, ok := params.Special[string(v)]
					if !ok {
						return "", fmt.Errorf("float special is not matching[%s]:[%s]", string(v), string(b))
					}
					if idx, ok := getSpecialIdx(tem, key.Name); ok {
						num, err := strconv.ParseFloat(string(op[idx]), 64)
						if err != nil {
							return "", fmt.Errorf("data[%s] can not convert to float64", string(op[idx]))
						}
						str = fmt.Sprintf(floatFormat, keyName, num)
					} else {
						str = fmt.Sprintf(floatFormat, keyName, 0.0)
					}
				} else {
					num, err := strconv.ParseFloat(string(op[ik]), 64)
					if err != nil {
						return "", fmt.Errorf("data[%s] can not convert to float64", string(op[ik]))
					}
					str = fmt.Sprintf(floatFormat, keyName, num)
				}
			default:
				if opless {
					tem, ok := params.Special[string(v)]
					if !ok {
						return "", fmt.Errorf("string special is not matching")
					}
					if idx, ok := getSpecialIdx(tem, key.Name); ok {
						str = fmt.Sprintf(stringFormat, keyName, string(op[idx]))
					} else {
						str = fmt.Sprintf(stringFormat, keyName, "")
					}
				} else {
					str = fmt.Sprintf(stringFormat, keyName, string(op[ik]))
				}
			}
			buf.WriteString(str)
		}

		bufWriteByte(buf, '}', inline)
	}
	bufWriteByte(buf, ']', inline)
	return buf.String(), nil
}

func inlineData(b []byte, k *Key) {

}
