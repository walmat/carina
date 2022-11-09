package taskutil

import "regexp"

func TrimLeadingZeroes(input string) string {
	return regexp.MustCompile(`^0+`).ReplaceAllString(input, "")
}
