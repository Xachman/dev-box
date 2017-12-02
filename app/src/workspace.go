package main

import (
	"bytes"
	"fmt"
	"os/exec"
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
		args := []string{}
		args = append(args, "run")
		args = append(args, "-d")
		args = append(args, "--name")
		args = append(args, fmt.Sprintf("%s_%s", config.Namespace, w.Name))

		for _, value := range w.Ports {
			args = append(args, "-p")
			args = append(args, fmt.Sprintf("0:%d", value))
		}

		for key, value := range w.Environment {
			args = append(args, "-e")
			args = append(args, fmt.Sprintf("%s=%s", key, value))
		}

		args = append(args, "-v")
		args = append(args, fmt.Sprintf("%s/%s:%s", config.GetVolumeDir(), w.Name, w.Volume))
		args = append(args, fmt.Sprintf("%s", w.Image))

		w.runCommand("docker", args)
		return
	}

	w.runCommand("docker", []string{"start", w.containerName()})
}
func (w *Workspace) remove() {
	w.runCommand("docker", []string{"rm", "-f", w.containerName()})
}
func (w *Workspace) Stop() {
	w.runCommand("docker", []string{"stop", w.containerName()})
}
func (w *Workspace) Status() string {
	res, err := w.runCommand("docker", []string{"inspect", "-f", "{{.State.Status}}", w.containerName()})
	if err != nil {
		return "not created"
	}
	return strings.TrimSpace(res)
}

func (w *Workspace) containerName() string {
	config := GetConfig()
	return config.GetNamespace() + "_" + w.Name
}
func (w *Workspace) getContainerId() string {
	//docker ps -aqf "name=
	res, _ := w.runCommand("docker", []string{"ps", "-aqf", "name=" + w.containerName()})
	return strings.TrimSpace(res)
}
func (w *Workspace) runCommand(url string, cmdArgs []string, method string) (string, error) {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
	}

	switch(method) {
	case "post":
		httpc.post("http://unix"+url, "application/json", params)
	}
	
}


func (w *Workspace) exists() bool {
	name, _ := w.runCommand("docker", []string{"ps", "-a", "-f", "name=" + w.containerName(), "--format", "{{.Names}}"})
	if w.containerName() == strings.TrimSpace(name) {
		return true
	}
	return false
}

type WorkspaceStatus struct {
	Status string
}
