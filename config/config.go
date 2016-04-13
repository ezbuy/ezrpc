package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ezbuy/statsd"
)

type config struct {
	Statsd *statsd.Config `json:"statsd"`
}

// Config is global config
var Config *config

// InitConfig load config from json file
func InitConfig(path string) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	Config = new(config)
	if err := json.Unmarshal(buf, Config); err != nil {
		return err
	}

	// init dependencies
	statsd.Setup(Config.Statsd)

	return nil
}
