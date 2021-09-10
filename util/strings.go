package util

import "strings"

func ReplaceAndTrim(s string) string {
	replacer := strings.NewReplacer("\n", "", "\t", "")
	return replacer.Replace(strings.TrimSpace(s))
}
