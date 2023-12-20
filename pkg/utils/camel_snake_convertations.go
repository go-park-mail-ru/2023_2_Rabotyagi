package utils

import (
	"regexp"
	"strings"
)

var (
	link          = regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func CamelToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToLower(snake)
}

func SnakeToCamelCase(str string) string {
	return link.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.ReplaceAll(s, "_", ""))
	})
}
