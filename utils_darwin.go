//go:build darwin

package main

import (
	"os"
)

const (
	LineEnding = "\n"
)

func GetHomeDir() string {
	return os.Getenv("HOME")
}

func EtcHostsFilename() string {
	return "/etc/hosts"
}
