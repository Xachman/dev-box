package main

import (
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	workspaceController := WorkspaceController{
		DataDir: "../data/workspaces",
	}
	consoleController := ConsoleController{}

	http.Handle("/workspaces/exec/", websocket.Handler(consoleController.ExecContainer))
	http.HandleFunc("/workspaces/start/", workspaceController.StartContainer)
	http.HandleFunc("/workspaces/stop/", workspaceController.StopContainer)
	http.HandleFunc("/workspaces/status/", workspaceController.ContainerStatus)
	http.HandleFunc("/workspaces/remove/", workspaceController.RemoveContainer)
	http.HandleFunc("/workspaces", workspaceController.Index)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(":9080", nil)
}
