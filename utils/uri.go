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

func GetUriQuery(uri string) string {
	query := ""
	splits := strings.Split(uri[1:], "/")
	if len(splits) >= 2 {
		query = splits[1]
	}
	return query
}
