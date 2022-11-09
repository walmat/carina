package proxies

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"os"
)

type exportedState struct {
	Groups  map[string]ProxyGroup `json:"groups"`
	ToGroup map[string]string     `json:"toGroup"`
	Links   map[string]string     `json:"links"`
}

func ExportAllGob(outPath string) error {
	var outBuf bytes.Buffer
	enc := gob.NewEncoder(&outBuf)

	outState := exportedState{
		Groups:  proxyGroups,
		ToGroup: proxyToGroup,
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

	proxiesMu.Lock()
	defer proxiesMu.Unlock()
	linksMu.Lock()
	defer linksMu.Unlock()

	for a, b := range state.Groups {
		proxyGroups[a] = b
	}

	for a, b := range state.ToGroup {
		proxyToGroup[a] = b
	}

	for a, b := range state.Links {
		links[a] = b
	}

	return nil
}

func ExportAllJson(outPath string) error {
	outState := exportedState{
		Groups:  proxyGroups,
		ToGroup: proxyToGroup,
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

	proxiesMu.Lock()
	defer proxiesMu.Unlock()
	linksMu.Lock()
	defer linksMu.Unlock()

	for a, b := range state.Groups {
		proxyGroups[a] = b
	}

	for a, b := range state.ToGroup {
		proxyToGroup[a] = b
	}

	for a, b := range state.Links {
		links[a] = b
	}

	return nil
}
