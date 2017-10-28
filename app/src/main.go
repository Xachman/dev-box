package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
)

func handler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("../data/workspaces")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {

		fmt.Fprintf(w, strings.Replace(f.Name(), ".yml", "", -1))
	}
}
func main() {
	workspaceController := WorkspaceController{
		DataDir: "../data/workspaces",
	}
	consoleController := ConsoleController{}

	http.Handle("/workspaces/exec/", websocket.Handler(consoleController.ExecContainer))
	http.HandleFunc("/workspaces/start/", workspaceController.StartContainer)
	http.HandleFunc("/workspaces/stop/", workspaceController.StopContainer)
	http.HandleFunc("/workspaces/status/", workspaceController.ContainerStatus)
	http.HandleFunc("/workspaces", workspaceController.Index)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(":9080", nil)
}
