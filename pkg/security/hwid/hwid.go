package hwid

import (
	"fmt"
	"nebula/pkg/security"
	"strconv"
)

var (
	cachedHwid security.HWID
)

func init() {
	cachedHwid = get()
}

func get() security.HWID {
	var cpuCores string
	if coreCount, err := GetCPUCores(); err != nil {
		cpuCores = fmt.Sprint("error:", err)
	} else {
		cpuCores = strconv.FormatUint(uint64(coreCount), 10)
	}

	var gpuName string
	if name, err := GetGPUName(); err != nil {
		gpuName = fmt.Sprint("error:", err)
	} else {
		gpuName = name
	}

	var diskSerials []string
	if serials, err := GetDiskSerials(); err != nil {
		diskSerials = []string{fmt.Sprint("error:", err)}
	} else {
		diskSerials = serials
	}

	var physMem string
	if mem, err := GetPhysicalMemory(); err != nil {
		physMem = fmt.Sprint("error:", err)
	} else {
		physMem = strconv.FormatInt(mem, 10)
	}

	var macAddr string
	if mac, err := GetMacAddr(); err != nil {
		macAddr = fmt.Sprint("error:", err)
	} else {
		macAddr = mac
	}

	return security.HWID{
		CpuName:  GetCPUName(),
		CpuCores: cpuCores,

		GpuName: gpuName,

		DiskSerials: diskSerials,

		PhysicalMemory: physMem,

		MacAddress: macAddr,
	}
}

func Get() security.HWID {
	return cachedHwid
}
