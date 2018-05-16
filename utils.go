package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"
)

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

func execCommand(command string) ([]byte, error) {
	cmd := exec.Command(command)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	slurp, _ := ioutil.ReadAll(stderr)
	if err := cmd.Wait(); err != nil {
		return slurp, nil
	}
	return slurp, nil
}
