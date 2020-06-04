package bitcoinwalletwatcher

import (
	"encoding/json"
	"io/ioutil"
)

// Config struct
type Config struct {
	InfoFile       string   `json:"info_filepath"`
	DefaultLoopSec float64  `json:"default_loop_in_sec"`
	Adresses       []string `json:"addresses"`
}

// NewConfig reads the path and returns watcher config
func NewConfig(path string) (*Config, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(f, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
