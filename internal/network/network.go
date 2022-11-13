package network

import (
	"errors"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Lofanmi/go-switch-hosts/contracts"
	"github.com/Lofanmi/go-switch-hosts/internal/gateway"
	"github.com/Lofanmi/go-switch-hosts/internal/gotil"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Network struct {
	ifaceList     contracts.NetInterfaceSlice
	handleManager contracts.PCAPHandleManager
}

func NewNetwork(handleManager contracts.PCAPHandleManager) (network contracts.Network, err error) {
	ifaceList, err := net.Interfaces()
	if err != nil {
		return
	}
	network = &Network{
		ifaceList:     ifaceList,
		handleManager: handleManager,
	}
	return
}

func (s *Network) Interfaces() contracts.NetInterfaceSlice {
	return s.ifaceList
}

func (s *Network) InterfaceIP(iface *net.Interface, to4 bool) net.IP {
	ips, err := iface.Addrs()
	if err != nil {
		return nil
	}
	for _, a := range ips {
		res := a.(*net.IPNet).IP
		if to4 && res.To4() != nil {
			return res
		}
		if !to4 && res.To16() != nil {
			return res
		}
	}
	return nil
}

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

func (s *Network) GatewayHardwareAddr(gateway net.IP, iface *net.Interface) (addr net.HardwareAddr, err error) {
	ip := s.InterfaceIP(iface, gateway.To4() != nil)
	start := time.Now()
	eth := layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(iface.HardwareAddr),
		SourceProtAddress: ip.To4(),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		DstProtAddress:    gateway.To4(),
	}
	buf := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	if err = gopacket.SerializeLayers(buf, options, &eth, &arp); err != nil {
		return
	}
	var handle *pcap.Handle
	if handle, err = s.handleManager.GetHandle(iface.Name); err != nil {
		return
	}
	if err = handle.WritePacketData(buf.Bytes()); err != nil {
		return
	}
	for {
		if time.Since(start) > time.Second*3 {
			err = errors.New("timeout getting ARP reply")
			return
		}
		var data []byte
		data, _, err = handle.ReadPacketData()
		if err == pcap.NextErrorTimeoutExpired {
			continue
		}
		if err != nil {
			return
		}
		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)
		if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
			a := arpLayer.(*layers.ARP)
			if a.Operation == layers.ARPReply && net.IP(a.SourceProtAddress).Equal(gateway) {
				return a.SourceHwAddress, nil
			}
		}
	}
}

func (s *Network) GetTCPConnectionList() (data contracts.TCPConnectionSlice, err error) {
	cmd := exec.Command("lsof", "-bnlPiTCP", "-i", "TCP", "-sTCP:ESTABLISHED")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	str := strings.ReplaceAll(string(output), "\r", "")
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}
		ii := strings.Index(line, "TCP")
		if ii == -1 {
			continue
		}
		ii += len("TCP ")
		jj := ii + strings.Index(line[ii:], " ")
		split := strings.Split(line[ii:jj], "->")
		src := strings.Split(split[0], ":")
		dst := strings.Split(split[1], ":")
		srcPort, _ := strconv.ParseInt(src[1], 10, 64)
		dstPort, _ := strconv.ParseInt(dst[1], 10, 64)
		data = append(data, contracts.TCPConnection{
			SrcIP:   net.ParseIP(src[0]),
			SrcPort: uint16(srcPort),
			DstIP:   net.ParseIP(dst[0]),
			DstPort: uint16(dstPort),
		})
	}
	return
}

func (s *Network) KillTCPConnection(connection contracts.TCPConnection) (err error) {
	return
}
