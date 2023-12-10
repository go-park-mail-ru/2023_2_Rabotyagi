package utils

import "strings"

func SimplifyPath(path string) string {
	prefixCut := "/api/v1/img"
	if strings.HasPrefix(path, prefixCut) {
		return prefixCut
	}

	return path
}
