package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	CfKey   string   `yaml:"cfKey"`
	CfEmail string   `yaml:"cfEmail"`
	Zone    string   `yaml:"zone"`
	Names   []string `yaml:"names"`
}

func LoadFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	config := &Config{}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, config)
	return config, err
}
