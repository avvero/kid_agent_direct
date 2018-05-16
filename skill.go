package main

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Skills struct {
	Version string  `yaml:"version"`
	Skills  []Skill `yaml:"skills"`
}

type Skill struct {
	Pattern  string   `yaml:"pattern"`
	Template string   `yaml:"template"`
	Scripts  []string `yaml:"script"`
}

func ReadSkills(fileName string) (*Skills, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	skills := Skills{}
	err = yaml.Unmarshal(content, &skills)
	if err != nil {
		return nil, err
	}
	log.Printf("skills: %s", skills)
	return &skills, nil
}
