package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
)
type ContainerController struct {
}

type ContainerResponse struct {
	Image 		string `json:"image"`
	Volume 		string `json:"volume"`
	VolumeDir 	string `json:"volumeDir"`
	Name 		string `json:"name"`
	Ports 		[]int  `json:"ports"`
	Env			map[string]string `json:"env"`
}
func (c *ContainerController) getContainer(image, name, volume, volumeDir string, ports []int, env map[string]string) Container {
	return NewContainer(image, volume, name, volumeDir, ports, env, Docker{})
}
func (c *ContainerController) Start(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(r.Body)
		fmt.Println(string(body))
		cRes := ContainerResponse{}
		json.Unmarshal(body, &cRes)
		container := c.getContainer(
			cRes.Image, 
			cRes.Name, 
			cRes.Volume, 
			cRes.VolumeDir, 
			cRes.Ports, 
			cRes.Env,
		)
		fmt.Println(cRes)
		container.start()
	} else {
		fmt.Fprintf(w, "Bad Method")
	}
}

func (c *ContainerController) Index() {
}
