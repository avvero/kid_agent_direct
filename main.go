package main

import (
	"encoding/json"
	"flag"
	"log"
	"runtime"
	"time"
)

var (
	pollEndpoint = flag.String("pollEndpoint", "", "poll endpoint")
	pullInterval = flag.Int("pullInterval", 1, "update interval for infos")
)

func main() {
	flag.Parse()

	log.Printf("Pulling from: %s", *pollEndpoint)

	_, err := ReadSkills("skills.yaml")
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
					handleTask(task)
				}
			}
		}
	}()
	runtime.Goexit()
}

func handleTask(task *Task) {
	log.Printf("Task is: %s", task.Value)
	stdout, err := execCommand(task.Value)
	if err != nil {
		log.Printf("Error during command execution: %s", err)
	}
	if stdout != nil {
		log.Printf("%s\n", stdout)
	}
}
