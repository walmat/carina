package hwid

import (
	"errors"
	"github.com/jaypipes/ghw"
	"strings"
)

func GetMacAddr() (string, error) {
	network, err := ghw.Network()
	if err != nil {
		return "", err
	}

	for _, nic := range network.NICs {
		if strings.Contains(nic.Name, "Wi-Fi") || strings.HasPrefix(nic.Name, "wlp") {
			return nic.MacAddress, nil
		}
	}
	return "", errors.New("not found")
}
