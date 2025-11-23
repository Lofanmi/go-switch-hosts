//go:build windows

package main

import (
	"os"
	"path/filepath"
)

const (
	LineEnding = "\r\n"
)

func GetHomeDir() string {
	home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return home
}

func EtcHostsFilename() string {
	return filepath.Join(os.Getenv("SYSTEMROOT"), "system32", "drivers", "etc", "hosts")
}
