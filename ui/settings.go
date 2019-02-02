package ui

import (
	"encoding/gob"
	"os"
	"path"
)

// Setting - interface type for one setting of the application
type Setting interface {
	Load(settings *Settings)
	Save(settings *Settings)
}

// Settings - manages all settings of the application
type Settings struct {
	VolumeLevel uint8
	settingList []Setting
}

// NewSettings - create a new Settings object
func NewSettings() *Settings {
	settings := Settings{}
	return &settings
}

// Load - load settings from a file
func (settings *Settings) Load() {
	filename := settingsPath()
	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
		decoder := gob.NewDecoder(file)
		if err := decoder.Decode(&settings); err != nil {
			settings.makeDefaultSettings()
		}
	} else {
		settings.makeDefaultSettings()
	}

	for _, v := range settings.settingList {
		v.Load(settings)
	}
}

// Save - save settings to a file
func (settings *Settings) Save() error {
	for _, v := range settings.settingList {
		v.Save(settings)
	}

	filename := settingsPath()
	dir, _ := path.Split(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	return encoder.Encode(settings)
}

// Register - add a new setting object to the list
func (settings *Settings) Register(setting Setting) {
	for _, v := range settings.settingList {
		if v == setting {
			return
		}
	}
	settings.settingList = append(settings.settingList, setting)
}

// Unregister - removes a setting object from the list
func (settings *Settings) Unregister(setting Setting) {
	for i, v := range settings.settingList {
		if v == setting {
			settings.settingList = append(settings.settingList[:i], settings.settingList[i+1:]...)
			return
		}
	}
}

func (settings *Settings) makeDefaultSettings() {
	settings.VolumeLevel = MaxVolumeLevel
}
