package network

import (
	"net"
	"os/exec"
	"strings"

	"github.com/Lofanmi/go-switch-hosts/internal/gateway"
	"github.com/Lofanmi/go-switch-hosts/internal/gotil"
)

func (s *Network) Route(ip net.IP) (gatewayIP net.IP, iface *net.Interface, err error) {
	cmd := exec.Command("route", "-n", "get", ip.String())
	data, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	str := strings.TrimSpace(gotil.StringCut(string(data), "gateway: ", "\n", false))
	if str == "" {
		return gateway.Default()
	}
	gatewayIP = net.ParseIP(str)
	ifaceName := strings.TrimSpace(gotil.StringCut(string(data), "interface: ", "\n", false))
	iface = s.ifaceList.ByName(ifaceName)
	return
}
