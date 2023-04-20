package parser2

import (
	"bytes"
	"sync"
)

var bufpool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer([]byte{})
	},
}

var standardTaFields = map[string]struct{}{
	"#account_id":     {},
	"#app_id":         {},
	"#distinct_id":    {},
	"#uuid":           {},
	"#type":           {},
	"#time":           {},
	"#ip":             {},
	"#event_name":     {},
	"#event_id":       {},
	"#first_check_id": {},
}
