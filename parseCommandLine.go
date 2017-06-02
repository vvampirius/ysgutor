package ysgutor

import (
	"regexp"
)

func ParseCommandLine(commandLine string) []string {
	parsedArray := make([]string, 0)
	re, _ := regexp.Compile(`([^\\]"(\\"|[^"])+"|((\\\s)|\S)+)`)
	for _, m := range re.FindAllStringSubmatch(commandLine, -1) {
		parsedArray = append(parsedArray, unQuote(m[0]))
	}
	return parsedArray
}

func unQuote(s string) string {
	if match := regexp.MustCompile(`^\s*"(.*)"$`).FindAllStringSubmatch(s, -1); len(match)==1 && len(match[0])==2 {
		s = match[0][1]
	}
	s = regexp.MustCompile(`\\"`).ReplaceAllString(s, `"`)
	return s
}
