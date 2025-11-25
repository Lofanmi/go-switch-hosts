//go:build darwin

package main

import (
	"os"

	"github.com/ying32/govcl/vcl/types"
)

const (
	LineEnding = "\n"

	CommandKeyCode uint32 = types.SsSuper
)

func GetHomeDir() string {
	return os.Getenv("HOME")
}

func EtcHostsFilename() string {
	return "/etc/hosts"
}
