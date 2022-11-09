package harvesters

func ImportState(harvestersPath string) error {
	if err := ImportHarvesters(harvestersPath); err != nil {
		return err
	}
	return nil
}

func ExportState(harvesterPath string) error {
	if err := ExportHarvesters(harvesterPath); err != nil {
		return err
	}

	return nil
}

func ImportGmailState(gmailsPath string) error {
	if err := ImportGmails(gmailsPath); err != nil {
		return err
	}
	return nil
}

func ExportGmailState(gmailsPath string) error {
	if err := ExportGmails(gmailsPath); err != nil {
		return err
	}

	return nil
}
