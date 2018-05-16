package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"regexp"
	"runtime"
	"time"
	"fmt"
)

var (
	pollEndpoint = flag.String("pollEndpoint", "", "poll endpoint")
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
		var v1 interface{}
		var v2 interface{}
		var v3 interface{}
		var v4 interface{}
		var v5 interface{}
		fmt.Sscanf(task.Value, matchedSkill.Template, &v1, &v2, &v3, &v4, &v5)
		log.Printf("%s", v1)
		log.Printf("%s", v2)
	
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
