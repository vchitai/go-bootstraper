package utils

import (
	"strings"
)

type NormalizedString []string

func ToSnake2(input string) string {
	//s := []rune(input)
	var res = ""
	//n := len(input)
	start := 0

	dfas := getDFAs()
	for _, c := range input {
		for _, dfa := range dfas {
			if gap := dfa.check(c); gap != nil {
				res += strings.ToLower(input[start:gap.PrevEnd]) + "_"
				start = gap.NextStart
			}
		}
	}
	res += strings.ToLower(input[start:])
	return res
}
func Normalize(s string) NormalizedString {
	res := make([]string, 0)
	for _, tokens := range strings.Split(s, "_") {
		for _, token := range strings.Split(tokens, "-") {
			res = append(res, strings.ToLower(token))
		}
	}
	return res
}

func (s NormalizedString) ToSnake() string {
	return strings.Join(s, "_")
}

func (s NormalizedString) WithDash() string {
	return strings.Join(s, "-")
}

func (s NormalizedString) ToCamel() string {
	if len(s) == 0 {
		return ""
	}
	res := s[0]
	for _, token := range s[1:] {
		if len(token) == 0 {
			continue
		}
		res += strings.ToUpper(token[:1]) + token[1:]
	}
	return res
}

func (s NormalizedString) ToTitle() string {
	res := ""
	for _, token := range s {
		if len(token) == 0 {
			continue
		}
		res += strings.ToUpper(token[:1]) + token[1:]
	}
	return res
}
