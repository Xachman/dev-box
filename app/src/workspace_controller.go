package main

import (
	"fmt"
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

func (ws *WorkspaceController) StartContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		c := Container{
			Image: "test",
			Path:  "/test",
		}

		c.Start()
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
				fmt.Fprintf(w, "%s <a href=\"\"", ws.Name)
			}
		}
	} else {
		fmt.Fprintf(w, "Bad Method")
	}
}
