package hwid

import (
	"github.com/jaypipes/ghw"
)

func GetDiskSerials() ([]string, error) {
	block, err := ghw.Block()
	if err != nil {
		return nil, err
	}

	var serials []string
	for _, disk := range block.Disks {
		serials = append(serials, disk.SerialNumber)
	}
	return serials, nil
}
