package utils

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"
	"bytes"
	"encoding/json"
	"text/template"
)

func HttpGet(url string) ([]byte, error) {
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

func HttpPost(url string, body interface{}) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	jsonStr, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
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

func ExecCommand(command string) ([]byte, error) {
	cmd := exec.Command("sh", "-c", command)
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

func ProcessTemplate(tmp string, keys map[string]string)(string, error)  {
	commandTemplate, err := template.New("command").Parse(tmp)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = commandTemplate.Execute(buf, keys)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}