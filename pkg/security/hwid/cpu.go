package hwid

import (
	"github.com/intel-go/cpuid"
	"github.com/jaypipes/ghw"
)

func GetCPUName() string {
	return cpuid.ProcessorBrandString
}

func GetCPUCores() (uint32, error) {
	cpu, err := ghw.CPU()
	if err != nil {
		return 0, err
	}

	return cpu.TotalCores, nil
}
