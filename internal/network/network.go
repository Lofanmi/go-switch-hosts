package network

import (
	"errors"
	"math/rand"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Lofanmi/go-switch-hosts/contracts"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	log "github.com/sirupsen/logrus"
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
	readTimeout := time.Second
	if err = s.killTCPConnection(connection.DstIP, connection.DstPort, readTimeout); err != nil {
		return
	}
	return s.killTCPConnection(connection.SrcIP, connection.SrcPort, readTimeout)
}

func (s *Network) killTCPConnection(dstIP net.IP, dstPort uint16, readTimeout time.Duration) (err error) {
	gw, iface, err := s.Route(dstIP)
	if err != nil {
		return
	}
	addr, err := s.GatewayHardwareAddr(gw, iface)
	if err != nil {
		return err
	}
	srcIP := s.InterfaceIP(iface, dstIP.To4() != nil)
	version, ethernetType := uint8(4), layers.EthernetTypeIPv4
	if dstIP.To4() == nil {
		version, ethernetType = uint8(6), layers.EthernetTypeIPv6
	}
	limit := 1024
	port := rand.Intn(65536-limit) + limit + 1
	eth := layers.Ethernet{SrcMAC: iface.HardwareAddr, DstMAC: addr, EthernetType: ethernetType}
	ip4 := layers.IPv4{SrcIP: srcIP, DstIP: dstIP, Version: version, TTL: 64, Protocol: layers.IPProtocolTCP}
	tcp := layers.TCP{SrcPort: layers.TCPPort(port), DstPort: layers.TCPPort(dstPort), Seq: rand.Uint32(), SYN: true, Window: 65535}
	if err = tcp.SetNetworkLayerForChecksum(&ip4); err != nil {
		return
	}
	buf := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	if err = gopacket.SerializeLayers(buf, options, &eth, &ip4, &tcp); err != nil {
		return
	}
	var handle *pcap.Handle
	if handle, err = s.handleManager.GetHandle(iface.Name); err != nil {
		return
	}
	start, ack, window := time.Now(), uint32(0), uint16(0)
	for {
		if time.Since(start) > readTimeout {
			err = errors.New("timeout getting ARP reply")
			return
		}
		if err = gopacket.SerializeLayers(buf, options, &eth, &ip4, &tcp); err != nil {
			return
		}
		if err = handle.WritePacketData(buf.Bytes()); err != nil {
			return
		}
		var data []byte
		if data, _, err = handle.ReadPacketData(); err != nil {
			if err != pcap.NextErrorTimeoutExpired {
				log.Println(err)
			}
			continue
		}
		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)
		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			continue
		}
		receiveTCP4 := packet.TransportLayer().(*layers.TCP)
		ack = receiveTCP4.Ack
		window = receiveTCP4.Window
		break
	}
	tcp = layers.TCP{SrcPort: layers.TCPPort(port), DstPort: layers.TCPPort(dstPort), Seq: ack, RST: true, SYN: true, Window: window}
	if err = tcp.SetNetworkLayerForChecksum(&ip4); err != nil {
		return
	}
	buf = gopacket.NewSerializeBuffer()
	if err = gopacket.SerializeLayers(buf, options, &eth, &ip4, &tcp); err != nil {
		return
	}
	return handle.WritePacketData(buf.Bytes())
}
