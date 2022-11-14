//go:build !windows

package gotil

import (
	"os"
)

func GetHomeDir() string {
	return os.Getenv("HOME")
}
