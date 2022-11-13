package gateway

import (
	"net"
	"os/exec"
	"strings"
)

func Default() (gateway net.IP, iface *net.Interface, err error) {
	cmd := exec.Command("netstat", "-nr")
	data, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	for _, line := range strings.Split(strings.Trim(string(data), " \r\n\t"), "\n") {
		line = strings.ReplaceAll(strings.Trim(line, " \t"), "\t", " ")
		if strings.HasPrefix(line, "default") {
			pieces := strings.Fields(line)
			if len(pieces) <= 1 || pieces[0] == "" {
				return
			}
			gateway = net.ParseIP(pieces[1])
			iface, err = net.InterfaceByName(pieces[3])
			return
		}
	}
	return
}
