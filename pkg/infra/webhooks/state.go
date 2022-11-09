package webhooks

func ImportState(webhooksPath string) error {
	if err := ImportWebhooks(webhooksPath); err != nil {
		return err
	}
	return nil
}

func ExportState(webhooksPath string) error {
	if err := ExportWebhooks(webhooksPath); err != nil {
		return err
	}

	return nil
}
