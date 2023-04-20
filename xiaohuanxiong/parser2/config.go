package parser2

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

var vip *viper.Viper

var cfg *config

type config struct {
	Mapping   map[string]string   `json:"mapping"`
	Templates map[string]*setting `json:"templates"`
}

type setting struct {
	Template  string               `json:"template"`
	Configs   []*propertiesSetting `json:"configs"`
	templates []string
	Filename  string
	idxs      map[string]int
}

type propertiesSetting struct {
	Type    string                 `json:"type"` // 事件类型
	Filter  map[string]interface{} `json:"filter"`
	Repeat  map[string]string      `json:"repeat"`
	Mapping map[string]string      `json:"mapping"`
	Fields  map[string]*field      `json:"fields"`
}

type field struct {
	Idx    int                  `json:"idx"`
	Type   string               `json:"type"`
	Params *StringToOtherParams `json:"params"` // 特殊处理参数，只应用于特殊字段
}

func InitConfig() error {
	if err := loadConfig(); err != nil {
		return err
	}
	return parseConfig()
}

func parseConfig() error {
	if err := vip.Unmarshal(&cfg); err != nil {
		return err
	}
	if cfg.Mapping == nil {
		return fmt.Errorf("mapping is not config")
	}
	if cfg.Templates == nil {
		return fmt.Errorf("templates is not config")
	}
	for f, setting := range cfg.Templates {
		setting.Filename = f
		setting.parseTemplate()
		for _, c := range setting.Configs {
			c.parseMapping()
		}
	}
	return nil
}

func (s *propertiesSetting) parseMapping() {
	if s.Mapping == nil {
		s.Mapping = cfg.Mapping
	}
}

func (s *setting) parseTemplate() {
	templates := strings.Split(s.Template, string([]byte{0x01}))
	s.templates = templates
	s.idxs = map[string]int{}
	for i, v := range s.templates {
		s.idxs[v] = i
	}
}

func loadConfig() error {
	vip = viper.NewWithOptions(viper.KeyDelimiter("::"))
	vip.SetConfigName("template")
	vip.AddConfigPath("./conf")
	return vip.ReadInConfig()
}
