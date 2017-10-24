package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
				fmt.Println(file)
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
				t.Execute(w, ws)
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
