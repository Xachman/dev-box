package main

import (
	"bytes";
	"context";
	"encoding/json";
	"io/ioutil";
	"net";
	"net/http";
	"fmt";
)

type Docker struct {

}
func (c Docker) call(url string, args string, method string) (*APIResponse, error) {
	///containers/create
	///containers/(id or name)/start
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
	}
	switch method {
	case "post":
		params := bytes.NewBufferString(args)
		resp, err := httpc.Post("http://unix"+url, "application/json", params)
		if err != nil {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			apiR := new(APIResponse)
			json.Unmarshal(body, &apiR)
			return apiR, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		apiR := new(APIResponse)
		json.Unmarshal(body, &apiR)
		return apiR, nil
	case "get":
		resp, err := httpc.Get("http://unix" + url)
		if err != nil {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			apiR := new(APIResponse)
			json.Unmarshal(body, &apiR)
			return apiR, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		apiR := new(APIResponse)
		json.Unmarshal(body, &apiR)
		return apiR, err
	case "delete":
		req, err := http.NewRequest("DELETE", "http://unix"+url, nil)
		if err != nil {
			panic(err)
		}
		resp, err := httpc.Do(req)
		if err != nil {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			apiR := new(APIResponse)
			json.Unmarshal(body, &apiR)
			fmt.Println(apiR.Message)
			return apiR, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		apiR := new(APIResponse)
		json.Unmarshal(body, &apiR)
		return apiR, err
	}
	return nil, nil
}