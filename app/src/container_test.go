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
	type MockCallable struct {}
	func (c MockCallable) call(url string, args string, method string) (*APIResponse, error) {
		fmt.Println(url)
		if url != "dude" {
			t.Errorf()
		}
		return &APIResponse{}, nil
	}
	mock := MockCallable{}
	container := NewContainer("basic/image", "/src", "myname","/dest", []int{80, 8080}, stringMap, mock)
    // if total != 10 {
    //    t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
	// }
	container.start()
	
}

