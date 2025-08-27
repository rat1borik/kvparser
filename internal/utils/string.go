package utils

import "regexp"

func IsEmptyOrWhitespace(s string) bool {
	return s == "" || !regexp.MustCompile(`\S+`).MatchString(s)
}
