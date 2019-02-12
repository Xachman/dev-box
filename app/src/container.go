package main

import (
	"fmt"
	"strings"
)

// Container starts and stops containers
type Container struct {
	Image         string
	Ports         []int
	Volume        string
	Name          string
	VolumeDir     string
	Environment   map[string]string
	ContainerName string
	Caller 		  Callable
}

// NewContainer Creates container
func NewContainer(img, vol, name, volDir string, ports []int, env map[string]string, caller Callable) Container {
	return Container{
		Image:       img,
		Volume:      vol,
		Name:        name,
		VolumeDir:   volDir,
		Ports:       ports,
		Environment: env,
		Caller: caller,
	}
}
func (c *Container) start() {
	if !c.exists() {
		c.pullImage()
		config := GetConfig()
		args := "{"

		args += "\"HostConfig\": {"
		for _, value := range c.Ports {
			args += fmt.Sprintf("\"PortBindings\": { \"%d/tcp\": [{ \"HostPort\": \"0\" }] },", value)
		}
		volumeDir := c.Name
		if c.VolumeDir != "" {
			volumeDir = c.VolumeDir
		}
		args += "\"Binds\": ["
		args += fmt.Sprintf("\"%s/%s:%s\"", config.GetVolumeDir(), volumeDir, c.Volume)
		args += "]"
		args += "},"

		args += "\"Env\": ["
		index := 0
		for key, value := range c.Environment {
			if index > 0 {
				args += ","
			}
			index++
			args += fmt.Sprintf("\"%s=%s\"", key, value)
		}
		args += "],"

		args += fmt.Sprintf("\"Image\": \"%s\"", c.Image)
		args += "}"
		fmt.Println(args)
		apiR, _ := c.runCommand(fmt.Sprintf("/containers/create?name=%s_%s", config.Namespace, c.Name), args, "post")
		fmt.Println("Message: " + apiR.Message)
		c.startContainer()
		return
	}

	c.startContainer()
}

func (c *Container) runCommand(url string, args string, method string) (*APIResponse, error) {
	return c.Caller.call(url, args, method)
}
func (c *Container) exists() bool {
	apiR, _ := c.runCommand(fmt.Sprintf("/containers/%s/json", c.containerName()), "", "get")
	fmt.Println("Message Exists: " + apiR.Message)
	if apiR != nil && strings.TrimSpace(apiR.Message) == "" {
		return true
	}
	return false
}

func (c *Container) pullImage() {
	c.runCommand(fmt.Sprintf("/images/create?fromImage=%s", c.Image), "", "post")
}

func (c *Container) startContainer() {
	apiR, _ := c.runCommand(fmt.Sprintf("/containers/%s/start", c.containerName()), "", "post")
	fmt.Println("Message: " + apiR.Message)
}

func (c *Container) containerName() string {
	config := GetConfig()
	return config.GetNamespace() + "_" + c.Name
}
func (c *Container) portmaps() string {
	apiR, _ := c.runCommand(fmt.Sprintf("/containers/%s/json", c.containerName()), "", "get")
	key := fmt.Sprintf("%d/tcp", c.Ports[0])
	fmt.Println(key)
	if val, ok := apiR.NetworkSettings.Ports[key]; ok {
		return val[0].HostPort
	}
	return ""
}
func (c *Container) remove() {
	c.stop()
	c.runCommand(fmt.Sprintf("/containers/%s", c.containerName()), "", "delete")
}

func (c *Container) stop() {
	c.runCommand(fmt.Sprintf("/containers/%s/stop", c.containerName()), "", "post")
}

func (c *Container) status() string {
	apiR, _ := c.runCommand(fmt.Sprintf("/containers/%s/json", c.containerName()), "", "get")
	fmt.Println(apiR.State)
	if apiR != nil && strings.TrimSpace(apiR.State.Status) != "" {
		return apiR.State.Status
	}
	return "not created"
}

func (c *Container) getContainerId() string {
	apiR, _ := c.runCommand(fmt.Sprintf("/containers/%s/json", c.containerName()), "", "get")
	if apiR != nil {
		return apiR.Id
	}
	return "false"
}

type APIResponse struct {
	Message string `json:"message"`
	State   struct {
		Status string
	}
	NetworkSettings struct {
		Ports map[string][]struct {
			HostPort string
		}
	}
	Id string
}
