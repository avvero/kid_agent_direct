package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"regexp"
	"runtime"
	"time"
)

var (
	pollEndpoint = flag.String("pollEndpoint", "https://f2g.site/bot/kid/api/tasks/29:1rxhBdD4tY9pijFOBxI4JatuCjaCxMFkKFgRszgxFbQ0/poll", "poll endpoint")
	pullInterval = flag.Int("pullInterval", 1, "update interval for infos")
)

func main() {
	flag.Parse()

	log.Printf("Pulling from: %s", *pollEndpoint)

	config, err := ReadConfiguration("config.yaml")
	if err != nil {
		log.Fatal("Error during skills parsing: %s", err)
	}

	ticker := time.NewTicker(time.Duration(*pullInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				body, err := callEndpoint(*pollEndpoint)
				if err != nil {
					log.Printf("Error during task pulling: %s", err)
				} else if len(body) > 0 {
					task := &Task{}
					json.Unmarshal(body, task)
					err := handleTask(config, task)
					log.Printf("Error during task handling: %s", err)
				}
			}
		}
	}()
	runtime.Goexit()
}

func handleTask(config *Configuration, task *Task) error {
	log.Printf("Task is: %s", task.Value)

	var matchedSkill *Skill

	for _, skill := range config.Skills {
		matched, err := regexp.MatchString(skill.Pattern, task.Value)
		if err != nil {
			return err
		}
		if matched {
			matchedSkill = &skill
			break
		}
	}
	if matchedSkill == nil {
		return errors.New("Can't handle task - don't know how")
	}
	if matchedSkill.Template == "" {
		stdout, err := execCommand(task.Value)
		if err != nil {
			return err
		}
		if stdout != nil {
			log.Printf("%s\n", stdout)
		}
		return nil
	} else {
		//lex := newLexer(task.Value)
		//go lex.tokenize()

		//log.Printf("%s\n", lex.tokens)

		for _, script := range matchedSkill.Scripts {
			stdout, err := execCommand(script)
			if err != nil {
				return err
			}
			if stdout != nil {
				log.Printf("%s\n", stdout)
			}
		}
		return nil
	}
}
