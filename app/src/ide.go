package main

import "fmt"

/*

mount workspace files

start ide container

*/

// Workspace as defined
type IDE struct {
	Image     string
	Ports     []int
	Volume    string
	Workspace Workspace
}

func (ide *IDE) Start() {
	c := ide.getContainer()
	c.start()
}

func (ide *IDE) getContainer() Container {
	environment := make(map[string]string)
	fmt.Println(ide.Workspace.Volume)
	return NewContainer(ide.Image, ide.Volume, "cloud9_"+ide.Workspace.Name, ide.Workspace.Name, ide.Ports, environment)
}

func (ide *IDE) stop() {
	c := ide.getContainer()
	c.stop()
}
func (ide *IDE) ports() {
	c := ide.getContainer()
	return c.portmaps()
}

func (ide *IDE) remove() {
	c := ide.getContainer()
	c.remove()
}
