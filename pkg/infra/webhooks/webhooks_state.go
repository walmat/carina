package webhooks

import (
	"encoding/json"
	"os"
)

func ImportWebhooks(webhooksPath string) error {
	webhookBytes, err := os.ReadFile(webhooksPath)
	if err != nil {
		return err
	}

	parsedWebhooks := make(map[string]WebhookData)
	if err = json.Unmarshal(webhookBytes, &parsedWebhooks); err != nil {
		return err
	}

	webhooksMutex.Lock()
	defer webhooksMutex.Unlock()
	webhooks = parsedWebhooks
	return nil
}

func ExportWebhooks(outPath string) error {
	webhooksMutex.RLock()
	defer webhooksMutex.RUnlock()
	exportedBytes, err := json.Marshal(webhooks)
	if err != nil {
		return err
	}

	return os.WriteFile(outPath, exportedBytes, os.ModePerm)
}
