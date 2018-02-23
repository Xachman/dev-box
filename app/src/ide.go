package main

/*

mount workspace files

start ide container

*/

// Workspace as defined
type IDE struct {
	Container
	Image     string
	Ports     []int
	Volume    string
	Workspace Workspace
}

func (ide *IDE) setValues() {
	environment := make(map[string]string)
	ide.Container = Container{
		Image:       ide.Image,
		Volume:      ide.Volume,
		Ports:       ide.Ports,
		Name:        "cloud9_" + ide.Workspace.Name,
		Environment: environment,
		VolumeDir:   ide.Workspace.Name,
	}
}
