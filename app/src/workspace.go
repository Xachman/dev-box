package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

// Workspace as defined
type Workspace struct {
	Image       string
	Ports       []int
	Volume      string
	Name        string
	VolumeDir   string
	Environment map[string]string
}

// Start starts Workspaces
func (w *Workspace) Start() {
	if !w.exists() {
		config := GetConfig()
		args := "{"

		args += "\"HostConfig\": {"
		for _, value := range w.Ports {
			args += fmt.Sprintf("\"PortBindings\": { \"%d/tcp\": [{ \"HostPort\": \"0\" }] },", value)
		}
		args += "\"Binds\": ["
		args += fmt.Sprintf("\"%s/%s:%s\"", config.GetVolumeDir(), w.Name, w.Volume)
		args += "]"
		args += "},"

		args += "\"Env\": ["
		index := 0
		for key, value := range w.Environment {
			if index > 0 {
				args += ","
			}
			index++
			args += fmt.Sprintf("\"%s=%s\"", key, value)
		}
		args += "],"

		args += fmt.Sprintf("\"Image\": \"%s\"", w.Image)
		args += "}"
		fmt.Println(args)
		w.runCommand(fmt.Sprintf("/containers/create?name=%s_%s", config.Namespace, w.Name), args, "post")
		return
	}

	//w.runCommand("docker", []string{"start", w.containerName()})
}
func (w *Workspace) remove() {
	//w.runCommand("docker", []string{"rm", "-f", w.containerName()})
}
func (w *Workspace) Stop() {
	//w.runCommand("docker", []string{"stop", w.containerName()})
}
func (w *Workspace) Status() string {
	// res, err := w.runCommand("docker", []string{"inspect", "-f", "{{.State.Status}}", w.containerName()})
	// if err != nil {
	// 	return "not created"
	// }
	return "not created"
}

func (w *Workspace) containerName() string {
	config := GetConfig()
	return config.GetNamespace() + "_" + w.Name
}
func (w *Workspace) getContainerId() string {
	//docker ps -aqf "name=
	//res, _ := w.runCommand("docker", []string{"ps", "-aqf", "name=" + w.containerName()})
	return "false"
}
func (w *Workspace) runCommand(url string, args string, method string) (*APIResponse, error) {
	///containers/create
	///containers/(id or name)/start
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
	}
	params := bytes.NewBufferString(args)
	switch method {
	case "post":
		resp, err := httpc.Post("http://unix"+url, "application/json", params)
		if err != nil {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			apiR := new(APIResponse)
			json.Unmarshal(body, &apiR)
			return apiR, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		apiR := new(APIResponse)
		json.Unmarshal(body, &apiR)
		return apiR, nil
	}
	return nil, nil
}

func (w *Workspace) exists() bool {
	// name, _ := w.runCommand("docker", []string{"ps", "-a", "-f", "name=" + w.containerName(), "--format", "{{.Names}}"})
	// if w.containerName() == strings.TrimSpace(name) {
	// 	return true
	// }
	return false
}

type WorkspaceStatus struct {
	Status string
}

type APIResponse struct {
	Message string `json:"message"`
}
