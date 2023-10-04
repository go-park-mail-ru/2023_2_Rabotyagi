package handler

import "strings"

func getPathParam(path string) string {
	last := strings.LastIndex(path, "/")
	if last == -1 {
		return ""
	}

	return path[last+1:]
}
