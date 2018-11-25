package ui

type VolumeController struct {
	audio *Audio
	level uint8
}

const (
	minLevel = 0
	maxLevel = 100
)

func NewVolumeController(audio *Audio) *VolumeController {
	volume := VolumeController{audio, maxLevel}
	volume.apply()
	return &volume
}

func (volume *VolumeController) Up() {
	if volume.level < maxLevel {
		volume.level += 10
		volume.apply()
	}
}

func (volume *VolumeController) Down() {
	if volume.level > minLevel {
		volume.level -= 10
		volume.apply()
	}
}

func (volume *VolumeController) apply() {
	normedLevel := float32(volume.level) / maxLevel
	volume.audio.SetVolume(normedLevel * normedLevel)
}
