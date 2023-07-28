package Utils

import (
	"encoding/json"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/Struct"
	"os"
)

func LoadConf(filePath string) (*Struct.Conf, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var conf = new(Struct.Conf)
	err = json.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}
	return conf, err
}
