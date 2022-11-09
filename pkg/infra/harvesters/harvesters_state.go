package harvesters

import (
	"encoding/json"
	"os"
)

func ImportHarvesters(harvestersPath string) error {
	harvestersBytes, err := os.ReadFile(harvestersPath)
	if err != nil {
		return err
	}

	parsedHarvesters := make(map[string]HarvesterData)
	if err = json.Unmarshal(harvestersBytes, &parsedHarvesters); err != nil {
		return err
	}

	harvestersMutex.Lock()
	defer harvestersMutex.Unlock()

	for _, h := range parsedHarvesters {
		if h.Solver != nil {
			h.Solver = nil
			parsedHarvesters[h.ID] = h
		}
	}

	Harvesters = parsedHarvesters
	return nil
}

func ExportHarvesters(outPath string) error {
	harvestersMutex.RLock()
	defer harvestersMutex.RUnlock()
	exportedBytes, err := json.Marshal(Harvesters)
	if err != nil {
		return err
	}

	return os.WriteFile(outPath, exportedBytes, os.ModePerm)
}

func ImportGmails(gmailsPath string) error {
	harvestersBytes, err := os.ReadFile(gmailsPath)
	if err != nil {
		return err
	}

	parsedHarvesters := make(map[string]GmailData)
	if err = json.Unmarshal(harvestersBytes, &parsedHarvesters); err != nil {
		return err
	}

	GmailsMutex.Lock()
	defer GmailsMutex.Unlock()
	Gmails = parsedHarvesters
	return nil
}

func ExportGmails(outPath string) error {
	GmailsMutex.RLock()
	defer GmailsMutex.RUnlock()
	exportedBytes, err := json.Marshal(Gmails)
	if err != nil {
		return err
	}

	return os.WriteFile(outPath, exportedBytes, os.ModePerm)
}
