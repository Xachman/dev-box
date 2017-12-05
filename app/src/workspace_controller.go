package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/websocket"
	yaml "gopkg.in/yaml.v2"
)

var host = flag.String("host", config.GetHost()+":"+config.GetPort(), "Docker host")

type WorkspaceController struct {
	DataDir string
}

func (wsc *WorkspaceController) StartContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "POST" {
		cName := strings.TrimPrefix(r.URL.Path, "/workspaces/start/")
		ws := wsc.getWorkspace(cName)

		ws.Start()
	} else {
		fmt.Fprintf(w, "Bad Method")
	}
}
func (wsc *WorkspaceController) StopContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "POST" {
		cName := strings.TrimPrefix(r.URL.Path, "/workspaces/stop/")
		ws := wsc.getWorkspace(cName)

		ws.Stop()
	} else {
		fmt.Fprintf(w, "Bad Method")
	}
}

func (wsc *WorkspaceController) RemoveContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "POST" {
		cName := strings.TrimPrefix(r.URL.Path, "/workspaces/remove/")
		ws := wsc.getWorkspace(cName)

		ws.remove()
	} else {
		fmt.Fprintf(w, "Bad Method")
	}
}
func (wsc *WorkspaceController) ContainerStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "GET" {
		cName := strings.TrimPrefix(r.URL.Path, "/workspaces/status/")
		ws := wsc.getWorkspace(cName)

		wss := WorkspaceStatus{
			Status: ws.Status(),
		}
		json, err := json.Marshal(wss)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	} else {
		fmt.Fprintf(w, "Bad Method")
	}
}
func (wsc *WorkspaceController) PortMaps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "GET" {
		cName := strings.TrimPrefix(r.URL.Path, "/workspaces/ports/")
		ws := wsc.getWorkspace(cName)
		hostport := ws.portmaps()
		config := GetConfig()
		port := struct {
			PortDomain string
		}{
			fmt.Sprintf("%s:%s", config.GetHost(), hostport),
		}
		json, err := json.Marshal(port)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	} else {
		fmt.Fprintf(w, "Bad Method")
	}
}

func (wsc *WorkspaceController) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "GET" {
		fileList := []string{}
		err := filepath.Walk(wsc.DataDir, func(path string, f os.FileInfo, err error) error {
			fileList = append(fileList, path)
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

		wsm := []Workspace{}
		for _, file := range fileList {
			if strings.Contains(file, ".yml") {
				ws := Workspace{}
				fileContents, err := ioutil.ReadFile(file)
				if err != nil {
					log.Fatalf("error: %v", err)
				}
				yErr := yaml.Unmarshal(fileContents, &ws)
				if yErr != nil {
					log.Fatalf("error: %v", err)
				}
				wsm = append(wsm, ws)
			}

		}
		json, err := json.Marshal(wsm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	} else {
		fmt.Fprintf(w, "Bad Method")
	}
}

func (wsc *WorkspaceController) getWorkspace(workspaceName string) Workspace {
	ws := Workspace{}
	fileContents, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.yml", wsc.DataDir, workspaceName))
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	yErr := yaml.Unmarshal(fileContents, &ws)
	if yErr != nil {
		log.Fatalf("error: %v", err)
	}
	return ws
}

func (wsc *WorkspaceController) ExecContainer(w *websocket.Conn) {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
	}
	containerName := w.Request().URL.Path[len("/workspaces/exec/"):]
	ws := wsc.getWorkspace(containerName)
	container := ws.getContainerId()
	if container == "" {
		w.Write([]byte("Container does not exist"))
		return
	}
	fmt.Println("container ok")
	type stuff struct {
		Id string
	}
	var s stuff
	params := bytes.NewBufferString("{\"AttachStdin\":true,\"AttachStdout\":true,\"AttachStderr\":true,\"Tty\":true,\"Cmd\":[\"/bin/bash\"]}")
	resp, err := httpc.Post("http://unix/containers/"+container+"/exec", "application/json", params)
	if err != nil {
		panic(err)
	}
	fmt.Println("http://"+*host+"/containers/"+container+"/exec", resp)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(data), &s)
	if err := hijack("/var/run/docker.sock", "POST", "/exec/"+s.Id+"/start", true, w, w, w, nil, nil); err != nil {
		panic(err)
	}
	fmt.Println("Connection!")
	fmt.Println(ws)
	spew.Dump(ws)
}

func (wsc *WorkspaceController) getJson(url string, target interface{}) error {
	var client = &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
