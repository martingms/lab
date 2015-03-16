package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Hosts []*Host
}

func ReadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	conf := Config{}

	err = decoder.Decode(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func (conf *Config) String() string {
	b, err := json.Marshal(conf)
	if err != nil {
		return ""
	}
	return string(b)
}
