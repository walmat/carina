package taskutil

import "strings"

func DecryptAkamaiSensor(body string) []string {
	const DecryptionString = "hftdrsetrTYE%TSTHTFjthdrghfgdhfjthrd"

	var result []rune
	for i := range body {
		result = append(result, []rune(body)[i]^[]rune(DecryptionString)[i%len(DecryptionString)])
	}

	parts := strings.Split(string(result), "*")
	return parts
}
