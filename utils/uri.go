package utils

import "strings"

func GetUriMethod(uri string) string {
	end := strings.Index(uri, "?")
	if end == -1 {
		end = len(uri)
	}
	splits := strings.Split(uri[1:end], "/")
	return splits[len(splits)-1]
}
