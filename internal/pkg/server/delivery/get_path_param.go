package delivery

import "strings"

func GetPathParam(path string) string {
	last := strings.LastIndex(path, "/")
	if last == -1 {
		return ""
	}

	return path[last+1:]
}
