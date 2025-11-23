package main

import (
	"os"
	"strings"
)

const (
	SwitchHostsContentStart = `# --- SWITCHHOSTS_CONTENT_START ---`
)

func ParseSystemHosts(etcHostsFilename string) string {
	data, err := os.ReadFile(etcHostsFilename)
	if err != nil {
		return ""
	}
	content := strings.TrimSpace(string(data))
	if content == "" {
		return ""
	}
	i := strings.Index(content, SwitchHostsContentStart)
	if i == -1 {
		return content
	}
	return strings.TrimSpace(content[:i])
}
