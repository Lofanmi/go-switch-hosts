package gotil

import (
	"strings"
)

func StringCut(s, begin, end string, withBegin bool) string {
	beginPos := strings.Index(s, begin)
	if beginPos == -1 {
		return ""
	}
	s = s[beginPos+len(begin):]
	endPos := strings.Index(s, end)
	if endPos == -1 {
		return ""
	}
	result := s[:endPos]
	if withBegin {
		return begin + result
	} else {
		return result
	}
}

func InArray[T comparable](e T, items []T) bool {
	for _, item := range items {
		if e == item {
			return true
		}
	}
	return false
}
