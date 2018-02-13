package main

import "fmt"

/*

mount workspace files

start ide container

*/

// Workspace as defined
type IDE struct {
	*Container
	Workspace Workspace
}

func (ide *IDE) setValues() {
	//environment := make(map[string]string)
	fmt.Println(ide.Workspace.Volume)

	ide.Name = "cloud9_" + ide.Workspace.Name
	//ide.VolumeDir =   volD
	//	Environment: environment,
	//}
}
