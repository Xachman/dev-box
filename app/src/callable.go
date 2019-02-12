package main

type Callable interface {
	call(url string, args string, method string) (*APIResponse, error)
}