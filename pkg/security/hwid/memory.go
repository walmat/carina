package hwid

import "github.com/jaypipes/ghw"

func GetPhysicalMemory() (int64, error) {
	mem, err := ghw.Memory()
	if err != nil {
		return 0, err
	}

	return mem.TotalPhysicalBytes, nil
}
