package ui

import (
	"math"

	"github.com/go-gl/gl/v2.1/gl"
)

// VolumeController - manages a volume level (with indication)
type VolumeController struct {
	audio     *Audio
	level     uint8
	hideTimer float64
}

// Minimal and maximal volume level
const (
	MinVolumeLevel = 0
	MaxVolumeLevel = 100
)

const (
	defaultHideTime = 1
	defaultHideA    = 0.1
	defaultHideB    = 0.9
)

func drawRect(x float32, y float32, w float32, h float32) {
	gl.Begin(gl.QUADS)
	gl.Vertex2f(x, y)
	gl.Vertex2f(x+w, y)
	gl.Vertex2f(x+w, y+h)
	gl.Vertex2f(x, y+h)
	gl.End()
}

func setBlend(value bool) {
	if value {
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	} else {
		gl.Disable(gl.BLEND)
	}
}

func alphaFunc(t float64) float64 {
	if t < defaultHideA {
		return t / defaultHideA
	}
	if t > defaultHideB {
		return (defaultHideTime - t) / (defaultHideTime - defaultHideB)
	}
	return 1.0
}

// NewVolumeController - creates a VolumeController object
func NewVolumeController(audio *Audio, settings *Settings) *VolumeController {
	volume := VolumeController{audio, MaxVolumeLevel, 0}
	settings.Register(&volume)
	return &volume
}

// Load - load a volume level from settings
func (volume *VolumeController) Load(settings *Settings) {
	volume.level = settings.VolumeLevel
	volume.apply()
}

// Save - save volume level to settings
func (volume *VolumeController) Save(settings *Settings) {
	settings.VolumeLevel = volume.level
}

// Up - try to raise a volume level
func (volume *VolumeController) Up() {
	if volume.level < MaxVolumeLevel {
		volume.level += 10
		volume.apply()
	}
	volume.runHideTimer()
}

// Down - try to lower a volume level
func (volume *VolumeController) Down() {
	if volume.level > MinVolumeLevel {
		volume.level -= 10
		volume.apply()
	}
	volume.runHideTimer()
}

func (volume *VolumeController) runHideTimer() {
	// (-inf, 0]
	if volume.hideTimer <= 0 {
		volume.hideTimer = defaultHideTime
		return
	}

	// (0, A)
	if volume.hideTimer < defaultHideA {
		volume.hideTimer = defaultHideTime - volume.hideTimer
		return
	}

	// [A, B]
	if volume.hideTimer <= defaultHideB {
		volume.hideTimer = defaultHideB
		return
	}

	// [B, defaultHideTime]
	// nothing to change
}

func (volume *VolumeController) apply() {
	normedLevel := float32(volume.level) / MaxVolumeLevel
	volume.audio.SetVolume(normedLevel * normedLevel)
}

// Draw - draw visual volume level indication
func (volume *VolumeController) Draw(x float32, y float32, w float32, h float32) {
	if volume.hideTimer <= 0 {
		return
	}

	normedLevel := float32(volume.level) / MaxVolumeLevel

	setBlend(true)

	alpha := uint8(math.Round(150 * alphaFunc(volume.hideTimer)))

	gl.Color4ub(17, 245, 34, alpha)
	drawRect(x, y, w*normedLevel, h)

	gl.Color4ub(219, 29, 20, alpha)
	drawRect(x+w*normedLevel, y, w*(1-normedLevel), h)

	gl.Color4ub(255, 255, 255, 255)
	setBlend(false)
}

// Update - update visual volume level indication
func (volume *VolumeController) Update(dt float64) {
	if volume.hideTimer > 0 {
		volume.hideTimer -= dt
	}
}
