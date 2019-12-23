package ui

import (
	"image"

	"github.com/bmv0/nes/nes"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const padding = 0

// GameView - manages all game related stuff (emulator, volume level, settings etc.)
type GameView struct {
	director *Director
	console  *nes.Console
	settings *Settings
	volume   *VolumeController
	title    string
	hash     string
	texture  uint32
	record   bool
	frames   []image.Image
	active   bool
}

// NewGameView - create new GameView object
func NewGameView(director *Director, console *nes.Console, title, hash string) View {
	texture := createTexture()
	settings := Settings{}
	volume := NewVolumeController(director.audio, &settings)
	newTitle := title + "  /  Save Alt+[0-9]  /  Load [0-9]"
	return &GameView{director, console, &settings, volume, newTitle, hash, texture, false, nil, false}
}

// Enter - called before switching to a GameView
func (view *GameView) Enter() {
	gl.ClearColor(0, 0, 0, 1)
	view.director.SetTitle(view.title)
	view.console.SetAudioChannel(view.director.audio.channel)
	view.console.SetAudioSampleRate(view.director.audio.sampleRate)

	window := view.director.window
	window.SetKeyCallback(view.onKey)
	window.SetFocusCallback(view.onFocus)
	view.active = isWindowActive(window)

	// load settings
	view.settings.Load()
	// load state
	if err := view.console.LoadState(savePath(view.hash)); err == nil {
		return
	}

	view.console.Reset()

	// load sram
	cartridge := view.console.Cartridge
	if cartridge.Battery != 0 {
		if sram, err := readSRAM(sramPath(view.hash)); err == nil {
			cartridge.SRAM = sram
		}
	}
}

// Exit - called before switching to other view
func (view *GameView) Exit() {
	window := view.director.window
	view.active = false
	window.SetKeyCallback(nil)
	window.SetFocusCallback(nil)
	view.console.SetAudioChannel(nil)
	view.console.SetAudioSampleRate(0)
	// save sram
	cartridge := view.console.Cartridge
	if cartridge.Battery != 0 {
		writeSRAM(sramPath(view.hash), cartridge.SRAM)
	}
	// save state
	view.console.SaveState(savePath(view.hash))
	// save settings
	view.settings.Save()
}

// Update - update and draw a game state
func (view *GameView) Update(t, dt float64) {
	if dt > 1 {
		dt = 0
	}
	window := view.director.window
	console := view.console

	if view.active {
		if joystickReset(glfw.Joystick1) {
			view.director.ShowMenu()
		}
		if joystickReset(glfw.Joystick2) {
			view.director.ShowMenu()
		}
		if readKey(window, glfw.KeyEscape) {
			view.director.ShowMenu()
		}
		updateControllers(window, console)
		console.StepSeconds(dt)
	}

	gl.BindTexture(gl.TEXTURE_2D, view.texture)
	setTexture(console.Buffer())
	drawBuffer(view.director.window)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	if view.record {
		view.frames = append(view.frames, copyImage(console.Buffer()))
	}

	view.volume.Update(dt)
	view.volume.Draw(0.65, -0.95, 0.3, 0.05)
}

func (view *GameView) onFocus(window *glfw.Window, focused bool) {
	view.active = focused
}

func (view *GameView) onKey(window *glfw.Window,
	key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch key {
		case glfw.KeyP:
			screenshot(view.console.Buffer())
		case glfw.KeyR:
			view.console.Reset()
		case glfw.KeyTab:
			if view.record {
				view.record = false
				animation(view.frames)
				view.frames = nil
			} else {
				view.record = true
			}
		case glfw.KeyPageUp:
			view.volume.Up()
		case glfw.KeyPageDown:
			view.volume.Down()
		case glfw.Key0, glfw.Key1, glfw.Key2, glfw.Key3, glfw.Key4, glfw.Key5, glfw.Key6, glfw.Key7, glfw.Key8, glfw.Key9:
			keyNum := int(key) - int(glfw.Key0)
			if mods == glfw.ModAlt {
				view.console.SaveState(statePath(view.hash, keyNum))
			} else if mods == 0 {
				view.console.LoadState(statePath(view.hash, keyNum))
			}
		}
	} else if action == glfw.Repeat {
		switch key {
		case glfw.KeyPageUp:
			view.volume.Up()
		case glfw.KeyPageDown:
			view.volume.Down()
		}
	}
}

func drawBuffer(window *glfw.Window) {
	w, h := window.GetFramebufferSize()
	s1 := float32(w) / 256
	s2 := float32(h) / 240
	f := float32(1 - padding)
	var x, y float32
	if s1 >= s2 {
		x = f * s2 / s1
		y = f
	} else {
		x = f
		y = f * s1 / s2
	}
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-x, -y)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(x, -y)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-x, y)
	gl.End()
}

func updateControllers(window *glfw.Window, console *nes.Console) {
	turbo := console.PPU.Frame%6 < 3
	k1 := readKeys(window, turbo)
	j1 := readJoystick(glfw.Joystick1, turbo)
	j2 := readJoystick(glfw.Joystick2, turbo)
	console.SetButtons1(combineButtons(k1, j1))
	console.SetButtons2(j2)
}
