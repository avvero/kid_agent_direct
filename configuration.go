package main

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	Version string  `yaml:"version"`
	Skills  []Skill `yaml:"skills"`
}

type Skill struct {
	Pattern string            `yaml:"pattern"`
	Tokens  map[string]string `yaml:"tokens"`
	Scripts []string          `yaml:"script"`
}

func ReadConfiguration(fileName string) (*Configuration, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	config := Configuration{}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	log.Printf("config: %s", config)
	return &config, nil
}
