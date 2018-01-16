package main

import (
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

var appPath = "/app"
var config Config

func main() {
	if len(os.Args) > 1 {
		appPath = os.Args[1]
	}

	config = GetConfig()
	workspaceController := WorkspaceController{
		DataDir: appPath + "/data/workspaces",
	}

	ideController := IDEController{
		DataDir:      appPath + "/data/ides",
		WorkspaceDir: appPath + "/data/workspaces",
	}

	http.Handle("/workspaces/exec/", websocket.Handler(workspaceController.ExecContainer))
	http.HandleFunc("/workspaces/start/", workspaceController.StartContainer)
	http.HandleFunc("/workspaces/stop/", workspaceController.StopContainer)
	http.HandleFunc("/workspaces/status/", workspaceController.ContainerStatus)
	http.HandleFunc("/workspaces/remove/", workspaceController.RemoveContainer)
	http.HandleFunc("/workspaces/ports/", workspaceController.PortMaps)
	http.HandleFunc("/workspaces", workspaceController.Index)

	http.HandleFunc("/ides/start/", ideController.startIDE)
	http.Handle("/", http.FileServer(http.Dir(appPath+"/public")))
	http.ListenAndServe(":9080", nil)
}
