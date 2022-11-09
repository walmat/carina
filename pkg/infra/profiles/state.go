package profiles

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"os"
)

type exportedState struct {
	Groups  map[string]ProfileGroup `json:"groups"`
	ToGroup map[string]string       `json:"toGroup"`
	Links   map[string]string       `json:"links"`
}

func ExportAllGob(outPath string) error {
	var outBuf bytes.Buffer
	enc := gob.NewEncoder(&outBuf)

	outState := exportedState{
		Groups:  profileGroups,
		ToGroup: profileToGroup,
		Links:   links,
	}
	if err := enc.Encode(outState); err != nil {
		return err
	}

	return os.WriteFile(outPath, outBuf.Bytes(), os.ModePerm)
}

func ImportGob(outPath string) error {
	data, err := os.ReadFile(outPath)
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(bytes.NewReader(data))

	var state exportedState
	if err = dec.Decode(&state); err != nil {
		return err
	}

	profilesMu.Lock()
	defer profilesMu.Unlock()
	linksMu.Lock()
	defer linksMu.Unlock()

	for a, b := range state.Groups {
		profileGroups[a] = b
	}

	for a, b := range state.ToGroup {
		profileToGroup[a] = b
	}

	for a, b := range state.Links {
		links[a] = b
	}

	return nil
}

func ExportAllJson(outPath string) error {
	outState := exportedState{
		Groups:  profileGroups,
		ToGroup: profileToGroup,
		Links:   links,
	}

	outBytes, err := json.Marshal(outState)
	if err != nil {
		return err
	}

	return os.WriteFile(outPath, outBytes, os.ModePerm)
}

func ImportJson(outPath string) error {
	data, err := os.ReadFile(outPath)
	if err != nil {
		return err
	}

	var state exportedState
	if err = json.Unmarshal(data, &state); err != nil {
		return err
	}

	profilesMu.Lock()
	defer profilesMu.Unlock()
	linksMu.Lock()
	defer linksMu.Unlock()

	for a, b := range state.Groups {
		profileGroups[a] = b
	}

	for a, b := range state.ToGroup {
		profileToGroup[a] = b
	}

	for a, b := range state.Links {
		links[a] = b
	}

	return nil
}
