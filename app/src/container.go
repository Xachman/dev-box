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

// Container starts and stops containers
type Container struct {
	Image       string
	Ports       []int
	Volume      string
	Name        string
	VolumeDir   string
	Environment map[string]string
}

// NewContainer Creates container
func NewContainer(img, vol, name, volDir string, ports []int, env map[string]string) Container {
	return Container{
		Image:       img,
		Volume:      vol,
		Name:        name,
		VolumeDir:   volDir,
		Ports:       ports,
		Environment: env,
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
		args += "\"Binds\": ["
		args += fmt.Sprintf("\"%s/%s:%s\"", config.GetVolumeDir(), c.Name, c.Volume)
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
func (c *Container) exists() bool {
	apiR, _ := c.runCommand(fmt.Sprintf("/containers/%s/json", c.containerName()), "", "get")
	if apiR != nil && strings.TrimSpace(apiR.Message) == "" {
		return true
	}
	return false
}
func (c *Container) runCommand(url string, args string, method string) (*APIResponse, error) {
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
			fmt.Println(apiR.Message)
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

func (c *Container) pullImage() {
	c.runCommand(fmt.Sprintf("/images/create?fromImage=%s", c.Image), "", "post")
}

func (c *Container) startContainer() {
	c.runCommand(fmt.Sprintf("/containers/%s/start", c.containerName()), "", "post")
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
