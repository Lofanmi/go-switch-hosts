package contracts

import (
	"fmt"
	"net"
	"strings"
)

type NetInterfaceSlice []net.Interface

func (s NetInterfaceSlice) ByName(name string) *net.Interface {
	for _, iface := range s {
		if iface.Name == name {
			return &iface
		}
	}
	return nil
}

type TCPConnection struct {
	SrcIP   net.IP
	SrcPort uint16
	DstIP   net.IP
	DstPort uint16
}

func (s TCPConnection) String() string {
	return fmt.Sprintf("%-15s:%d   --->   %-15s:%d", s.SrcIP, s.SrcPort, s.DstIP, s.DstPort)
}

type TCPConnectionSlice []TCPConnection

func (s TCPConnectionSlice) FindByDstIP(dstIP net.IP) *TCPConnection {
	for _, c := range s {
		if dstIP.Equal(c.DstIP) {
			return &c
		}
	}
	return nil
}

func (s TCPConnectionSlice) String() string {
	var b strings.Builder
	for _, connection := range s {
		b.WriteString(connection.String())
		b.WriteByte('\n')
	}
	return b.String()
}

type Network interface {
	Interfaces() NetInterfaceSlice
	InterfaceIP(iface *net.Interface, to4 bool) net.IP
	Route(ip net.IP) (gatewayIP net.IP, iface *net.Interface, err error)
	GatewayHardwareAddr(gateway net.IP, iface *net.Interface) (addr net.HardwareAddr, err error)
	GetTCPConnectionList() (data TCPConnectionSlice, err error)
	KillTCPConnection(connection TCPConnection) (err error)
}
