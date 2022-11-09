package security

type HWID struct {
	CpuName  string `json:"cpun"`
	CpuCores string `json:"cpuc"`

	GpuName string `json:"gpu"`

	DiskSerials []string `json:"disk"`

	PhysicalMemory string `json:"mem"`

	MacAddress string `json:"mac"`
}

func contains(slice []string, value string) bool {
	for _, val := range slice {
		if val == value {
			return true
		}
	}
	return false
}

func HwidVariance(hwid1 HWID, hwid2 HWID) int {
	variance := 0
	if hwid1.CpuName != hwid2.CpuName {
		variance++
	} else if hwid1.CpuCores != hwid2.CpuCores {
		variance++
	}

	if hwid1.GpuName != hwid2.GpuName {
		variance++
	}

	hasMatchingDisks := false
	if len(hwid1.DiskSerials) >= len(hwid2.DiskSerials) {
		for _, serial1 := range hwid1.DiskSerials {
			if contains(hwid2.DiskSerials, serial1) {
				hasMatchingDisks = true
				break
			}
		}
	} else {
		for _, serial2 := range hwid2.DiskSerials {
			if contains(hwid1.DiskSerials, serial2) {
				hasMatchingDisks = true
				break
			}
		}
	}
	if hasMatchingDisks {
		variance++
	}

	if hwid1.PhysicalMemory != hwid2.PhysicalMemory {
		variance++
	}

	if hwid1.MacAddress != hwid2.MacAddress {
		variance++
	}

	return variance
}
