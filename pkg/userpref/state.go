package userpref

import (
	"encoding/json"
	"os"
)

var (
	preferences = UserPreferences{
		Language: LabelValue{
			Label: "English (US)",
			Value: "en-US",
		},
		Currency: LabelValue{
			Label: "USD ($)",
			Value: "usd",
		},
		Sounds: Sounds{
			Checkout: Sound{
				Name:    "checkout.mp3",
				Data:    nil, // TODO
				Type:    "mp3",
				Volume:  100,
				Default: true,
			},
			Harvester: Sound{
				Name:    "harvester.mp3",
				Data:    nil, // TODO
				Type:    "mp3",
				Volume:  100,
				Default: true,
			},
		},
		Behaviors: Behaviors{
			Notifications:        true,
			RetryCheckouts:       false,
			MonitorPooling:       false,
			AutoClickHarvester:   true,
			AutoLaunchCollective: true,
		},
	}
)

func Get() UserPreferences {
	return preferences
}

func Set(p UserPreferences) {
	preferences = p
}

type exportedState struct {
	UserPreferences `json:"userPreferences"`
	Sounds          `json:"sounds"`
	ThemeState      uint8 `json:"themeState"`
}

func ImportPreferences(outPath string) error {
	taskBytes, err := os.ReadFile(outPath)
	if err != nil {
		return err
	}

	var state exportedState
	if err = json.Unmarshal(taskBytes, &state); err != nil {
		return err
	}

	preferences = state.UserPreferences
	themeState = state.ThemeState
	return nil
}

func ExportPreferences(outPath string) error {
	exportedBytes, err := json.Marshal(exportedState{
		UserPreferences: preferences,
		ThemeState:      themeState,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(outPath, exportedBytes, os.ModePerm)
}
