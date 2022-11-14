package utils

import (
	"regexp"
	"strings"
)

// CheckResourcePath checks whether a resource is formatted
func CheckResourcePath(path string) bool {
	sp := strings.Split(path, ".")
	pattern := "^[A-Za-z]*.[A-Za-z0-9_-]*.[A-Za-z0-9_-]*$"
	ok, _ := regexp.MatchString(pattern, path)
	if len(sp) != 3 || !ok {
		return false
	}
	return true
}
