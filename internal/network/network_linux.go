package network

func (s *Network) Route(ip net.IP) (gatewayIP net.IP, iface *net.Interface, err error) {
	cmd := exec.Command("ip", "route", "get", ip.String())
	data, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	gatewayIP := strings.TrimSpace(StringCut(string(data), "via ", " ", false))
	gateway = net.ParseIP(gatewayIP)
	ifaceName := strings.TrimSpace(StringCut(string(data), "dev ", " ", false))
	iface = ifaces.ByName(ifaceName)
	return
}
