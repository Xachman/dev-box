package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
func startContainer(w http.ResponseWriter, r *http.Request) {
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
func main() {
	http.HandleFunc("/", handler)
	workspaceController := WorkspaceController{
		DataDir: "../data/workspaces",
	}
	http.HandleFunc("/workspaces/start", workspaceController.StartContainer)
	http.HandleFunc("/workspaces", workspaceController.Index)
	http.ListenAndServe(":9080", nil)
}
