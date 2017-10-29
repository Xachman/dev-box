package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/websocket"

	"gopkg.in/yaml.v2"
)

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
	containerName := w.Request().URL.Path[len("/workspaces/exec/"):]
	ws := wsc.getWorkspace(containerName)
	container := ws.getContainerId()
	fmt.Println(container)
	if container == "" {
		w.Write([]byte("Container does not exist"))
		return
	}
	type stuff struct {
		Id string
	}
	var s stuff
	params := bytes.NewBufferString("{\"AttachStdin\":true,\"AttachStdout\":true,\"AttachStderr\":true,\"Tty\":true,\"Cmd\":[\"/bin/bash\"]}")
	resp, err := http.Post("http://"+*host+"/containers/"+container+"/exec", "application/json", params)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(data), &s)
	if err := hijack(*host, "POST", "/exec/"+s.Id+"/start", true, w, w, w, nil, nil); err != nil {
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
