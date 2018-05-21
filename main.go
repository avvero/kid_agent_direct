package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"regexp"
	"runtime"
	"time"
	"text/template"
	"bytes"
	"fmt"
)

var (
	pullInterval = flag.Int("pullInterval", 1, "update interval for infos")
)

func main() {
	flag.Parse()

	config, err := ReadConfiguration("config.yaml")
	if err != nil {
		log.Fatal("Error during skills parsing: %s", err)
	}

	pollEndpoint := fmt.Sprintf("https://f2g.site/bot/kid/api/tasks/%s/poll", config.Channel)
	log.Printf("Pulling from: %s", pollEndpoint)


	ticker := time.NewTicker(time.Duration(*pullInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				body, err := callEndpoint(pollEndpoint)
				if err != nil {
					log.Printf("Error during task pulling: %s", err)
				} else if len(body) > 0 {
					task := &Task{}
					json.Unmarshal(body, task)
					err := handleTask(config, task)
					if err != nil {
						//TODO should reply to the kid
						log.Printf("Error during task handling: %s", err)
					} 					
				}
			}
		}
	}()
	runtime.Goexit()
}

func handleTask(config *Configuration, task *Task) error {
	log.Println("--------")
	log.Printf("Task is: %s", task.Value)

	var matchedSkill *Skill

	for _, skill := range config.Skills {
		matched, err := regexp.MatchString(skill.Pattern, task.Value)
		if err != nil {
			return err
		}
		if matched {
			matchedSkill = skill
			break
		}
	}
	if matchedSkill == nil {
		//TODO should reply to the kid
		return errors.New("Can't handle task - don't know how")
	}
	log.Printf("Skill:  %v", matchedSkill)
	lex := newLexer(task.Value)
	lex.tokenize()
	log.Println("Tokens:")
	for _, v := range  lex.tokens {
		log.Printf("  %v", v)
	}
	// parse keys
	keys := make(map[string]string)
	keys["COMMAND"] = task.Value
	for v, skillToketRegex := range matchedSkill.TokensRegex {
		for _, stringTokens := range lex.tokens {
			if skillToketRegex.MatchString(stringTokens.value) {
				keys[v] = stringTokens.value
			}
		}
	}
	log.Println("Keys:")
	for k, v := range keys {
		log.Printf("  %s: %s", k, v)
	}

	for _, script := range matchedSkill.Scripts {
		commandTemplate, err := template.New("command").Parse(script)
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		err = commandTemplate.Execute(buf, keys)
		if err != nil {
			return err
		}

		//out, err := exec.Command("sh","-c",buf.String()).Output()
		log.Printf("Command: %s\n", buf.String())
		out, err := execCommand(buf.String())
		if err != nil {
			return err
		}
		//TODO should reply to the kid
		log.Printf("Command out: %s\n", out)
	}
	return nil
}
