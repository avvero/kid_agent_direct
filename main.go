package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
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

type Task struct {
	Value string `json:"value,omitempty"`
}

func handleTask(task *Task) {
	log.Printf("Task is: %s", task.Value)

	cmd := exec.Command(task.Value)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
	if err := cmd.Start(); err != nil {
		log.Printf("Error: %s", err)
		return
	}
	slurp, _ := ioutil.ReadAll(stderr)
	log.Printf("%s\n", slurp)
	if err := cmd.Wait(); err != nil {
		log.Printf("Error: %s", err)
	}
}

func callEndpoint(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Status " + resp.Status)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
