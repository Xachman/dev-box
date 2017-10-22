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

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":9080", nil)
}
