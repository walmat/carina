package security

import (
	"encoding/json"
	"nebula/pkg/security/hwid"
	"testing"
)

const (
	HWIDWindows1 = `{"cpun":"AMD Ryzen 9 5900X 12-Core Processor            \u0000","cpuc":"12","gpu":"NVIDIA GeForce RTX 3080","disk":["E823_8FA6_BF53_0001_001B_448B_49DC_392F.","0025_38BA_0150_62B0.","{9f2fbac0-9c4e-4a9e-ae6e-7926534c7d3b}"],"mem":"34359738368","mac":"DC:1B:A1:4D:94:A5"}`
	HWIDWindows2 = `{"cpun":"AMD Ryzen 9 5900X 12-Core Processor            \u0000","cpuc":"12","gpu":"NVIDIA GeForce RTX 3080","disk":["E823_8FA6_BF53_0001_001B_448B_49DC_392F.","0025_38BA_0150_62B0.","{9f2fbac0-9c4e-4a9e-ae6e-7926534c7d3b}"],"mem":"34359738368","mac":"DC:1B:A1:4D:94:A5"}`
)

func unmarshalHwid(data []byte) (HWID, error) {
	var r HWID
	err := json.Unmarshal(data, &r)
	return r, err
}

func TestHWIDVarianceSame(t *testing.T) {
	hwid1, err := unmarshalHwid([]byte(HWIDWindows1))
	if err != nil {
		t.Fatal("error unmarhsalling hwid:", err)
	}

	if variance := hwid.Variance(hwid1, hwid1); variance > 0 {
		t.Fatalf("expected 0 variance, got %d variance", variance)
	}
}

func TestHWIDVarianceVaried(t *testing.T) {
	hwid1, err := unmarshalHwid([]byte(HWIDWindows1))
	if err != nil {
		t.Fatal("error unmarhsalling hwid:", err)
	}
	hwid2 := hwid1

	hwid2.GpuName = "GeForce GTX 1050 Ti"
	if variance := hwid.Variance(hwid1, hwid2); variance != 1 {
		t.Fatalf("expected 1 variance, got %d variance", variance)
	}

	hwid2.CpuName = "Intel(R) Core(TM) i7-4770HQ CPU @ 2.20GHz"
	if variance := hwid.Variance(hwid1, hwid2); variance != 2 {
		t.Fatalf("expected 2 variance, got %d variance", variance)
	}

	hwid2.CpuCores = "0"
	if variance := hwid.Variance(hwid1, hwid2); variance != 2 {
		t.Fatalf("expected 2 variance, got %d variance", variance)
	}

	hwid2.PhysicalMemory = "1024"
	if variance := hwid.Variance(hwid1, hwid2); variance != 3 {
		t.Fatalf("expected 3 variance, got %d variance", variance)
	}

	hwid2.MacAddress = "7B:81:8F:90:F2:3A"
	if variance := hwid.Variance(hwid1, hwid2); variance != 4 {
		t.Fatalf("expected 4 variance, got %d variance", variance)
	}

	hwid2.DiskSerials = []string{"0025_38BA_0150_62B0."}
	if variance := hwid.Variance(hwid1, hwid2); variance != 6 {
		t.Fatalf("expected 6 variance, got %d variance", variance)
	}
}
