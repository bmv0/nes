package ui

import "github.com/gordonklaus/portaudio"

// Audio - manages game audio
type Audio struct {
	stream         *portaudio.Stream
	sampleRate     float64
	outputChannels int
	channel        chan float32
	volume         float32
}

// NewAudio - create new Audio object
func NewAudio() *Audio {
	a := Audio{}
	a.channel = make(chan float32, 44100)
	return &a
}

// Start - begin audio playing
func (a *Audio) Start() error {
	host, err := portaudio.DefaultHostApi()
	if err != nil {
		return err
	}
	parameters := portaudio.HighLatencyParameters(nil, host.DefaultOutputDevice)
	stream, err := portaudio.OpenStream(parameters, a.callback)
	if err != nil {
		return err
	}
	if err := stream.Start(); err != nil {
		return err
	}
	a.stream = stream
	a.sampleRate = parameters.SampleRate
	a.outputChannels = parameters.Output.Channels
	return nil
}

// Stop - end audio playing
func (a *Audio) Stop() error {
	return a.stream.Close()
}

// SetVolume - set new volume level
func (a *Audio) SetVolume(newVolume float32) {
	a.volume = newVolume
}

func (a *Audio) callback(out []float32) {
	var output float32
	for i := range out {
		if i%a.outputChannels == 0 {
			select {
			case sample := <-a.channel:
				output = sample
			default:
				output = 0
			}
		}
		out[i] = output * a.volume
	}
}
