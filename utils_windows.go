//go:build windows

package main

import (
	"os"
	"path/filepath"

	"github.com/ying32/govcl/vcl/types"
	_ "github.com/ying32/liblclbinres"
)

const (
	LineEnding = "\r\n"

	CommandKeyCode uint32 = types.SsCtrl
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
