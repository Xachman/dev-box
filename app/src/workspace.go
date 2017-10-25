package main

import (
	"bytes"
	"fmt"
	"os"
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
	config := GetConfig()
	args := []string{}
	args = append(args, "run")
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

	fmt.Println(strings.Join(args, ", "))
	cmd := exec.Command("docker", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		os.Exit(1)
	}
	fmt.Println(out.String())
}
