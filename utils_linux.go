//go:build linux

package main

import (
	"os"

	_ "github.com/ying32/liblclbinres"
)

const (
	LineEnding = "\n"

	CommandKeyCode uint32 = types.SsCtrl
)

func GetHomeDir() string {
	return os.Getenv("HOME")
}

func EtcHostsFilename() string {
	return "/etc/hosts"
}
