package ui

import (
	"log"

	"github.com/bmv0/nes/nes"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

// View - view abstract interface
type View interface {
	Enter()
	Exit()
	Update(t, dt float64)
}

// Director - manages a main window, audio and switches views
type Director struct {
	window    *glfw.Window
	audio     *Audio
	view      View
	menuView  View
	timestamp float64
}

// NewDirector - create a new Director object
func NewDirector(window *glfw.Window, audio *Audio) *Director {
	director := Director{}
	director.window = window
	director.audio = audio
	return &director
}

// SetTitle - set title for a main window
func (d *Director) SetTitle(title string) {
	d.window.SetTitle(title)
}

// SetView - switche to other view
func (d *Director) SetView(view View) {
	if d.view != nil {
		d.view.Exit()
	}
	d.view = view
	if d.view != nil {
		d.view.Enter()
	}
	d.timestamp = glfw.GetTime()
}

// Step - one step of a main loop
func (d *Director) Step() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	timestamp := glfw.GetTime()
	dt := timestamp - d.timestamp
	d.timestamp = timestamp
	if d.view != nil {
		d.view.Update(timestamp, dt)
	}
}

// Start - launch a main loop
func (d *Director) Start(paths []string) {
	d.menuView = NewMenuView(d, paths)
	if len(paths) == 1 {
		d.PlayGame(paths[0])
	} else {
		d.ShowMenu()
	}
	d.Run()
}

// Run - main loop
func (d *Director) Run() {
	for !d.window.ShouldClose() {
		d.Step()
		d.window.SwapBuffers()
		glfw.PollEvents()
	}
	d.SetView(nil)
}

// PlayGame - show a game view
func (d *Director) PlayGame(path string) {
	hash, err := hashFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	console, err := nes.NewConsole(path)
	if err != nil {
		log.Fatalln(err)
	}
	d.SetView(NewGameView(d, console, path, hash))
}

// ShowMenu - show a menu view
func (d *Director) ShowMenu() {
	d.SetView(d.menuView)
}
