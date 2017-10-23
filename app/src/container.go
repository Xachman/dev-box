package main

import "fmt"

// Container starts and stops containers
type Container struct {
	Image, Path string
}

func (c *Container) Start() {
	fmt.Printf("docker run --rm /workspacepath:%s %s", c.Path, c.Image)
}
