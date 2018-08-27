package regex

import (
	"strings"
)

func regexMatch(toMatch, regex string) bool {
	regexSplit := strings.Split(trimRegex(regex), "*")
	count := 0
	lenToMatch := len(toMatch)

	for _, part := range regexSplit {
		l := len(part)
		if l > 1 {
			if part[0:l-1] != toMatch[count:count+l-1] {
				return false
			}
			count = count + l - 1
		}
		for count < lenToMatch && toMatch[count] == part[l-1] {
			count++
		}
	}
	if count == len(toMatch) {
		return true
	}
	return false
}
func trimRegex(regex string) string {
	regexNew := ""
	prev := false
	for i := 0; i < len(regex); i++ {
		if string(regex[i]) == "*" {
			if prev {
				continue
			} else {
				prev = true
				regexNew += string(regex[i])
			}
		} else {
			regexNew += string(regex[i])
			prev = false
		}
	}
	return regexNew
}
