package hwid

import (
	"errors"
	"github.com/jaypipes/ghw"
)

func GetGPUName() (string, error) {
	gpus, err := ghw.GPU()
	if err != nil {
		return "", err
	}

	for _, gpu := range gpus.GraphicsCards {
		return gpu.DeviceInfo.Product.Name, nil
	}
	return "", errors.New("not found")
}
