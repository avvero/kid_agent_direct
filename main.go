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
	"os"
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
	log.Printf("--------", task.Value)
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
		//TODO should reply to the kid
		return errors.New("Can't handle task - don't know how")
	}
	dictionary := make(map[string]*regexp.Regexp)
	for k, v := range matchedSkill.Tokens {
		dictionary[k] = regexp.MustCompile(v)
	}
	lex := newLexer(dictionary, task.Value)
	lex.tokenize()
	keys := make(map[string]string)
	for _, v := range  lex.tokens {
		keys[v.tokenType] = v.value
		log.Printf("%v", v)
	}

	for _, script := range matchedSkill.Scripts {
		command, err := template.New("command").Parse(script)
		if err != nil { panic(err) }
		err = command.Execute(os.Stdout, keys)
		if err != nil { panic(err) }

		//stdout, err := execCommand(command)
		//if err != nil {
		//	return err
		//}
		//if stdout != nil {
		//	TODO should reply to the kid
			//log.Printf("%s\n", stdout)
		//}
	}
	return nil
}
