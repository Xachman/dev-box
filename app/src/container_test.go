package main

import (
	"testing";
	"fmt";
)

func TestContainerStart(t *testing.T) {
	stringMap := map[string]string {
		"var": "key",
		"var2": "key2",
	}
	mock := MockCallable{
		t: t,
		test: AssertContainerStart,
	}
	container := NewContainer("basic/image", "/src", "myname","/dest", []int{80, 8080}, stringMap, mock)
	container.startContainer()
	
}

func AssertContainerStart(url string, args string, method string, t *testing.T) {
	fmt.Println(url)
	config := GetConfig()
	if url != fmt.Sprintf("/containers/%s_myname/start", config.GetNamespace()) {
		t.Errorf("url expected dude got %s", url)
	}
}

type MockCallable struct {
	test func(url string, args string, method string, t *testing.T)
	t *testing.T
}

func (c MockCallable) call(url string, args string, method string) (*APIResponse, error) {
	c.test(url, args, method, c.t)
	return &APIResponse{}, nil
}
