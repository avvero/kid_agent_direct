package main

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"regexp"
	"log"
)

type Configuration struct {
	Version string   `yaml:"version"`
	Host    string   `yaml:"host"`
	Channel string   `yaml:"channel"`
	Skills  []*Skill `yaml:"skills"`
}

type Skill struct {
	Pattern     string            `yaml:"pattern"`
	Tokens      map[string]string `yaml:"tokens"`
	TokensRegex map[string]*regexp.Regexp
	Scripts     []string          `yaml:"script"`
	Message     *Message          `yaml:"message"`
}

type Message struct {
	Channel string `yaml:"channel"`
	Text    string `yaml:"text"`
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
	//todo too dirty
	for _, skill := range config.Skills {
		skill.TokensRegex = make(map[string]*regexp.Regexp)
		for k, v := range skill.Tokens {
			skill.TokensRegex[k] = regexp.MustCompile(v)
		}
	}
	log.Printf("config: %s", config)
	return &config, nil
}