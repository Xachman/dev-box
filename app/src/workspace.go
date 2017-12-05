package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
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
		apiR, _ := w.runCommand(fmt.Sprintf("/containers/create?name=%s_%s", config.Namespace, w.Name), args, "post")
		fmt.Println("Message: " + apiR.Message)
		w.startContainer()
		return
	}

	w.startContainer()
}
func (w *Workspace) remove() {
	w.Stop()
	w.runCommand(fmt.Sprintf("/containers/%s", w.containerName()), "", "delete")
}
func (w *Workspace) Stop() {
	w.runCommand(fmt.Sprintf("/containers/%s/stop", w.containerName()), "", "post")
}
func (w *Workspace) startContainer() {
	w.runCommand(fmt.Sprintf("/containers/%s/start", w.containerName()), "", "post")
}
func (w *Workspace) Status() string {
	apiR, _ := w.runCommand(fmt.Sprintf("/containers/%s/json", w.containerName()), "", "get")
	fmt.Println(apiR.State)
	if apiR != nil && strings.TrimSpace(apiR.State.Status) != "" {
		return apiR.State.Status
	}
	return "not created"
}

func (w *Workspace) containerName() string {
	config := GetConfig()
	return config.GetNamespace() + "_" + w.Name
}
func (w *Workspace) getContainerId() string {
	apiR, _ := w.runCommand(fmt.Sprintf("/containers/%s/json", w.containerName()), "", "get")
	if apiR != nil {
		return apiR.Id
	}
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
	switch method {
	case "post":
		params := bytes.NewBufferString(args)
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
	case "get":
		resp, err := httpc.Get("http://unix" + url)
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
		return apiR, err
	case "delete":
		req, err := http.NewRequest("DELETE", "http://unix"+url, nil)
		if err != nil {
			panic(err)
		}
		resp, err := httpc.Do(req)
		if err != nil {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			apiR := new(APIResponse)
			json.Unmarshal(body, &apiR)
			return apiR, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		apiR := new(APIResponse)
		json.Unmarshal(body, &apiR)
		return apiR, err
	}
	return nil, nil
}
func (w *Workspace) exists() bool {
	apiR, _ := w.runCommand(fmt.Sprintf("/containers/%s/json", w.containerName()), "", "get")
	if apiR != nil && strings.TrimSpace(apiR.Message) == "" {
		return true
	}
	return false
}

type WorkspaceStatus struct {
	Status string
}

type APIResponse struct {
	Message string `json:"message"`
	State   struct {
		Status string
	}
	Id string
}
