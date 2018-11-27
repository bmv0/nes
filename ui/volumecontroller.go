package ui

type VolumeController struct {
	audio *Audio
	level uint8
}

const (
	MinVolumeLevel = 0
	MaxVolumeLevel = 100
)

func NewVolumeController(audio *Audio, settings *Settings) *VolumeController {
	volume := VolumeController{audio, MaxVolumeLevel}
	settings.Register(&volume)
	return &volume
}

func (volume *VolumeController) Load(settings *Settings) {
	volume.level = settings.VolumeLevel
	volume.apply()
}

func (volume *VolumeController) Save(settings *Settings) {
	settings.VolumeLevel = volume.level
}

func (volume *VolumeController) Up() {
	if volume.level < MaxVolumeLevel {
		volume.level += 10
		volume.apply()
	}
}

func (volume *VolumeController) Down() {
	if volume.level > MinVolumeLevel {
		volume.level -= 10
		volume.apply()
	}
}

func (volume *VolumeController) apply() {
	normedLevel := float32(volume.level) / MaxVolumeLevel
	volume.audio.SetVolume(normedLevel * normedLevel)
}
