package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"runtime"
	"time"
	"fmt"
	"github.com/avvero/kid_agent_direct/api"
	"github.com/avvero/kid_agent_direct/utils"
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

	pollEndpoint := fmt.Sprintf("%s/api/tasks/%s/poll", config.Host, config.Channel)
	log.Printf("Pulling from: %s", pollEndpoint)

	apiClient := api.NewApiClient(config.Host)

	ticker := time.NewTicker(time.Duration(*pullInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				body, err := utils.HttpGet(pollEndpoint)
				if err != nil {
					log.Printf("Error during task pulling: %s", err)
				} else if len(body) > 0 {
					task := &Task{}
					json.Unmarshal(body, task)
					err := handleTask(config, apiClient, task)
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

func handleTask(config *Configuration, apiClient *api.ApiClient, task *Task) error {
	log.Println("--------")
	log.Printf("Task is: %s", task.Value)

	skill,_ := config.findSkill(task.Value)
	if skill == nil {
		//TODO should reply to the kid
		return errors.New("Can't handle task - don't know how")
	}
	log.Printf("Skill:  %v", skill)
	// Key extractions
	keys := getCommandKeys(skill, task.Value)
	// Command execution
	// Script goes first
	for _, script := range skill.Scripts {
		command, err:= utils.ProcessTemplate(script, keys)

		//out, err := exec.Command("sh","-c",buf.String()).Output()
		log.Printf("Command: %s\n", command)
		out, err := utils.ExecCommand(command)
		if err != nil {
			return err
		}
		//TODO should reply to the kid
		log.Printf("Command out: %s\n", out)
	}
	// Send message if needed
	if skill.Message != nil {
		message, err:= utils.ProcessTemplate(skill.Message.Text, keys)
		if err != nil {
			return err
		}
		log.Printf("Send message: %s\n", message)
		return apiClient.SendMessage(skill.Message.Channel, message)
	}
	return nil
}
