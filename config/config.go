package config

import (
	"encoding/json"
	"io/ioutil"
)

type Postgres struct {
	Host       string `yaml:"host" json:"host"`
	Port       int    `yaml:"port" json:"port"`
	User       string `yaml:"user" json:"user"`
	Database   string `yaml:"database" json:"database"`
	Password   string `yaml:"password" json:"password"`
	Timeout    int    `yaml:"timeout" json:"timeout"`
	RetriesNum int    `yaml:"retries_num" json:"retries_num"`
}

type Config struct {
	Port     int       `yaml:"port" json:"port"`
	Postgres *Postgres `yaml:"postgres" json:"postgres"`
}

func Parse(filePath string) (*Config, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	err = json.Unmarshal(b, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
