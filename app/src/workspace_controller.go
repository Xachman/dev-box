package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type WorkspaceController struct {
	DataDir string
}

func (wsc *WorkspaceController) StartContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		cName := strings.TrimPrefix(r.URL.Path, "/workspaces/start/")
		ws := wsc.getWorkspace(cName)

		ws.Start()
	} else {
		fmt.Fprintf(w, "Bad Method")
	}
}
func (wsc *WorkspaceController) StopContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		cName := strings.TrimPrefix(r.URL.Path, "/workspaces/stop/")
		ws := wsc.getWorkspace(cName)

		ws.Stop()
	} else {
		fmt.Fprintf(w, "Bad Method")
	}
}

func (wsc *WorkspaceController) ContainerStatus(w http.ResponseWriter, r *http.Request) {
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
	if r.Method == "GET" {
		fileList := []string{}
		err := filepath.Walk(wsc.DataDir, func(path string, f os.FileInfo, err error) error {
			fileList = append(fileList, path)
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

		for _, file := range fileList {
			templateFile, _ := ioutil.ReadFile("./tmpl/workspaces/index.html")
			if strings.Contains(file, ".yml") {
				wsid := strings.TrimSuffix(filepath.Base(file), ".yml")
				wss := WorkspaceStatus{}
				wsc.getJson("http://localhost:9080/workspaces/status/"+wsid, &wss)
				fmt.Println(wss.Status)
				ws := Workspace{}
				fileContents, err := ioutil.ReadFile(file)
				if err != nil {
					log.Fatalf("error: %v", err)
				}
				yErr := yaml.Unmarshal(fileContents, &ws)
				if yErr != nil {
					log.Fatalf("error: %v", err)
				}
				t, _ := template.New("workspace").Parse(string(templateFile[:]))
				type M map[string]string
				t.Execute(w, M{
					"Name":   ws.Name,
					"Status": wss.Status,
				})
			}
		}
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

func (wsc *WorkspaceController) getJson(url string, target interface{}) error {
	var client = &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
