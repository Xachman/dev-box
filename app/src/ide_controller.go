package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type IDEController struct {
	DataDir      string
	WorkspaceDir string
}

func (idec *IDEController) startIDE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "POST" {
		cName := strings.TrimPrefix(r.URL.Path, "/ides/start/")
		fmt.Println("starting ide " + cName)
		IDE := idec.getIDE(cName)

		IDE.Start()
	} else {
		fmt.Fprintf(w, "Bad Method")
	}

}

func (idec *IDEController) getIDE(workspaceName string) IDE {
	ide := IDE{}
	ws := Workspace{}
	fileContents, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.yml", idec.DataDir, "cloud9"))
	workspaceContents, err2 := ioutil.ReadFile(fmt.Sprintf("%s/%s.yml", idec.WorkspaceDir, workspaceName))
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if err2 != nil {
		log.Fatalf("error: %v", err)
	}
	yErr := yaml.Unmarshal(fileContents, &ide)
	if yErr != nil {
		log.Fatalf("error: %v", err)
	}
	yErr = yaml.Unmarshal(workspaceContents, &ws)
	if yErr != nil {
		log.Fatalf("error: %v", err)
	}
	ide.Workspace = ws
	return ide
}
