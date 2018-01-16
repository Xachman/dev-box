package main

/*

mount workspace files

start ide container

*/

// Workspace as defined
type IDE struct {
	Image     string
	Ports     []int
	Workspace Workspace
}

func (ide *IDE) Start() {
	c := ide.getContainer()
	c.start()
}

func (ide *IDE) getContainer() Container {
	environment := make(map[string]string)
	return NewContainer(ide.Image, ide.Workspace.Volume, "cloud9"+"_"+ide.Workspace.Name, ide.Workspace.VolumeDir, ide.Ports, environment)
}
