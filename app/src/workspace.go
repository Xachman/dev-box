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
	if !w.exists() {
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

		w.runCommand("docker", args)
		return
	}

	w.runCommand("docker", []string{"start", w.containerName()})
}

func (w *Workspace) Stop() {
	w.runCommand("docker", []string{"stop", w.containerName()})
}
func (w *Workspace) Status() string {
	return strings.TrimSpace(w.runCommand("docker", []string{"inspect", "-f", "{{.State.Status}}", w.containerName()}))
}

func (w *Workspace) containerName() string {
	config := GetConfig()
	return config.GetNamespace() + "_" + w.Name
}
func (w *Workspace) runCommand(cmdString string, cmdArgs []string) string {
	fmt.Printf("%s with args: %s\n", cmdString, cmdArgs)
	cmd := exec.Command("docker", cmdArgs...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		os.Exit(1)
	}
	wErr := cmd.Wait()
	if wErr != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		os.Exit(1)
	}

	return out.String()
}

func (w *Workspace) exists() bool {
	name := w.runCommand("docker", []string{"ps", "-a", "-f", "name=" + w.containerName(), "--format", "{{.Names}}"})
	if w.containerName() == strings.TrimSpace(name) {
		return true
	}
	return false
}

type WorkspaceStatus struct {
	Status string
}
