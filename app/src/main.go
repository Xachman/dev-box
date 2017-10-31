package main

import (
	"net/http"

	"golang.org/x/net/websocket"
)

var config = GetConfig()

func main() {
	workspaceController := WorkspaceController{
		DataDir: "/app/data/workspaces",
	}

	http.Handle("/workspaces/exec/", websocket.Handler(workspaceController.ExecContainer))
	http.HandleFunc("/workspaces/start/", workspaceController.StartContainer)
	http.HandleFunc("/workspaces/stop/", workspaceController.StopContainer)
	http.HandleFunc("/workspaces/status/", workspaceController.ContainerStatus)
	http.HandleFunc("/workspaces/remove/", workspaceController.RemoveContainer)
	http.HandleFunc("/workspaces", workspaceController.Index)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(":9080", nil)
}
