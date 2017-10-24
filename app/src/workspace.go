package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
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
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("run --name %s_%s ", config.Namespace, w.Name))

	for _, value := range w.Ports {
		buffer.WriteString(fmt.Sprintf("-p 0:%d ", value))
	}

	for key, value := range w.Environment {
		buffer.WriteString(fmt.Sprintf("-e %s=%s ", key, value))
	}

	buffer.WriteString(fmt.Sprintf("-v %s/%s:%s %s", config.GetVolumeDir(), w.Name, w.Volume, w.Image))

	fmt.Println(buffer.String())
	if err := exec.Command("cat", "../").Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
