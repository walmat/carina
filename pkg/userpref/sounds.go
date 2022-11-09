package userpref

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	speakerOnce sync.Once
)

type (
	Sound struct {
		Name    string `json:"name"`
		Data    []byte `json:"data"`
		Type    string `json:"type"`
		Volume  int    `json:"volume"`
		Default bool   `json:"default"`
	}

	Sounds struct {
		Checkout  Sound `json:"checkout"`
		Harvester Sound `json:"harvester"`
	}
)

func decodeSound(data []byte, fileExt string) (beep.StreamSeekCloser, beep.Format, error) {
	dataReader := io.NopCloser(bytes.NewReader(data))

	switch fileExt {
	case ".mp3":
		return mp3.Decode(dataReader)
	case ".wav":
		return wav.Decode(dataReader)
	case ".ogg":
		return vorbis.Decode(dataReader)
	case ".flac":
		return flac.Decode(dataReader)
	default:
		return nil, beep.Format{}, errors.New("unsupported file format")
	}
}

func SetCheckoutSound(path string) error {
	if path == "" {
		preferences.Sounds.Checkout.Data = nil
		preferences.Sounds.Checkout.Type = ".mp3"
		preferences.Sounds.Checkout.Default = true
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var splitPath []string
	if runtime.GOOS == "windows" {
		splitPath = strings.Split(path, "\\")
	} else {
		splitPath = strings.Split(path, "/")
	}

	fileName := splitPath[len(splitPath)-1]
	fileExt := filepath.Ext(fileName)

	streamer, _, err := decodeSound(data, fileExt)
	if err != nil {
		return err
	}
	streamer.Close()

	preferences.Sounds.Checkout = Sound{
		Name:    fileName,
		Data:    data,
		Type:    fileExt,
		Default: false,
	}
	return nil
}

func SetCheckoutVolume(volume int) {
	checkoutSound := preferences.Sounds.Checkout
	checkoutSound.Volume = volume
	preferences.Sounds.Checkout = checkoutSound
}

func PlayCheckout() error {
	checkoutSound := preferences.Sounds.Checkout

	// Exit early if volume is 0
	if checkoutSound.Volume == 0 {
		return nil
	}

	if checkoutSound.Default {
		_ = json.Unmarshal([]byte(DefaultCheckout), &checkoutSound.Data)
	}

	streamer, format, err := mp3.Decode(io.NopCloser(bytes.NewReader(checkoutSound.Data)))
	if err != nil {
		return err
	}
	defer streamer.Close()

	resampler := beep.ResampleRatio(32, 1, streamer)
	volume := &effects.Volume{Streamer: resampler, Base: 5, Volume: float64(checkoutSound.Volume) / 100.}

	speakerOnce.Do(func() {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/5))
	})
	speaker.Play(volume)
	return nil
}

func SetHarvesterSound(path string) error {
	if path == "" {
		preferences.Sounds.Harvester.Data = nil
		preferences.Sounds.Harvester.Type = ".mp3"
		preferences.Sounds.Harvester.Default = true
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var splitPath []string
	if runtime.GOOS == "windows" {
		splitPath = strings.Split(path, "\\")
	} else {
		splitPath = strings.Split(path, "/")
	}

	fileName := splitPath[len(splitPath)-1]
	fileExt := filepath.Ext(fileName)

	streamer, _, err := decodeSound(data, fileExt)
	if err != nil {
		return err
	}
	streamer.Close()

	preferences.Sounds.Harvester = Sound{
		Name:    fileName,
		Data:    data,
		Type:    fileExt,
		Default: false,
	}
	return nil
}

func SetHarvesterVolume(volume int) {
	harvesterSound := preferences.Sounds.Harvester
	harvesterSound.Volume = volume
	preferences.Sounds.Harvester = harvesterSound
}

func PlayHarvester() error {
	harvesterSound := preferences.Sounds.Harvester

	// Exit early if volume is 0
	if harvesterSound.Volume == 0 {
		return nil
	}

	if harvesterSound.Default {
		_ = json.Unmarshal([]byte(DefaultHarvester), &harvesterSound.Data)
	}

	streamer, format, err := mp3.Decode(io.NopCloser(bytes.NewReader(harvesterSound.Data)))
	if err != nil {
		return err
	}
	defer streamer.Close()

	resampler := beep.ResampleRatio(32, 1, streamer)
	volume := &effects.Volume{Streamer: resampler, Base: 5, Volume: float64(harvesterSound.Volume) / 100.}

	speakerOnce.Do(func() {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/5))
	})
	speaker.Play(volume)
	return nil
}

func StopAllSounds() {
	speaker.Clear()
}
