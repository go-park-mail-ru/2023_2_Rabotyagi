package utils

import "strings"

func SimplifyPath(path string) string {
	prefixCut := "/img/"
	if strings.HasPrefix(path, prefixCut) {
		return prefixCut
	}

	return path
}
