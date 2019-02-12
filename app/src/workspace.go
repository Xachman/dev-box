package main

// Workspace as defined
type Workspace struct {
	Image       string
	Ports       []int
	Volume      string
	Name        string
	VolumeDir   string
	Environment map[string]string
}

// Start starts Workspaces
func (w *Workspace) Start() {
	c := w.getContainer()
	c.start()
}
func (w *Workspace) getContainer() Container {
	return NewContainer(w.Image, w.Volume, w.Name, w.VolumeDir, w.Ports, w.Environment, Docker{})
}
func (w *Workspace) launchIde(ide string) {

}
func (w *Workspace) remove() {
	c := w.getContainer()
	c.remove()
}
func (w *Workspace) Stop() {
	c := w.getContainer()
	c.stop()
}

func (w *Workspace) portmaps() string {
	c := w.getContainer()
	return c.portmaps()
}

func (w *Workspace) Status() string {
	c := w.getContainer()
	return c.status()
}

func (w *Workspace) getContainerId() string {
	c := w.getContainer()
	return c.getContainerId()
}

type WorkspaceStatus struct {
	Status string
}
