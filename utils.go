package go_mongo

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strings"
)

func ConvertToCamelCase(s string) string {
	words := strings.Split(s, "_")
	for i := 1; i < len(words); i++ {
		words[i] = cases.Title(language.Und).String(words[i])
	}
	return strings.Join(words, "")
}

func ToSnakeCase(input string) string {
	// 匹配大写字母前面的位置
	regex := regexp.MustCompile(`([A-Z])`)
	snakeCase := regex.ReplaceAllString(input, "_$1")
	snakeCase = strings.ToLower(snakeCase)
	snakeCase = strings.TrimPrefix(snakeCase, "_")
	return snakeCase
}
