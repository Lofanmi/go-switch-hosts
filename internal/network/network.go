package network

import (
	"errors"
	"fmt"
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

func (s *Network) KillTCPConnection(connection *contracts.TCPConnection) (err error) {
	readTimeout := time.Second * 5
	log.Debugf("杀死连接 1 %s", connection)
	if err = s.killTCPConnection(connection.DstIP, connection.DstPort, connection.SrcPort, readTimeout); err != nil {
		log.Debugf(err.Error())
	}
	time.Sleep(time.Second)
	log.Debugf("杀死连接 2 %s", connection)
	if err = s.killTCPConnection(connection.DstIP, connection.DstPort, connection.SrcPort, readTimeout); err != nil {
		log.Debugf(err.Error())
	}
	return
}

func (s *Network) killTCPConnection(dstIP net.IP, dstPort, srcPort uint16, readTimeout time.Duration) (err error) {
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
	random := rand.Uint32()
	eth := layers.Ethernet{SrcMAC: iface.HardwareAddr, DstMAC: addr, EthernetType: ethernetType}
	var ip gopacket.SerializableLayer
	ip = &layers.IPv4{SrcIP: srcIP, DstIP: dstIP, Version: version, TTL: 64, Protocol: layers.IPProtocolTCP}
	if ethernetType == layers.EthernetTypeIPv6 {
		ip = &layers.IPv6{SrcIP: srcIP, DstIP: dstIP, Version: version} // IPv6 尚未测试
	}
	tcp := layers.TCP{SrcPort: layers.TCPPort(srcPort), DstPort: layers.TCPPort(dstPort), Seq: random, SYN: true, Window: 65535}
	if err = tcp.SetNetworkLayerForChecksum(ip.(gopacket.NetworkLayer)); err != nil {
		return
	}
	buf := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	if err = gopacket.SerializeLayers(buf, options, &eth, ip, &tcp); err != nil {
		return
	}
	var handle *pcap.Handle
	if handle, err = s.handleManager.GetHandle(iface.Name); err != nil {
		return
	}
	start, seq, ack, window := time.Now(), uint32(0), uint32(0), uint16(0)
	for {
		if time.Since(start) > readTimeout {
			err = fmt.Errorf("读取远端 [%s:%d] 超时，等待时间大于 [%s]", dstIP.String(), dstPort, readTimeout)
			return
		}
		if err = gopacket.SerializeLayers(buf, options, &eth, ip, &tcp); err != nil {
			return
		}
		log.Debugf("发送数据包 [SYN] %s", layersTCP2String(srcIP.String(), dstIP.String(), &tcp))
		if err = handle.WritePacketData(buf.Bytes()); err != nil {
			return
		}
		time.Sleep(time.Millisecond * 100)
		_, _, _ = handle.ReadPacketData() // 忽略第一个包，测试发现是上面的发包。
		// ignore, _, _ := handle.ReadPacketData() // 忽略第一个包，测试发现是上面的发包。
		// ignorePacket := gopacket.NewPacket(ignore, layers.LayerTypeEthernet, gopacket.NoCopy)
		// log.Println("------------------------------------------------------------------------------------------")
		// log.Println(ignorePacket.String())
		// log.Println("------------------------------------------------------------------------------------------")
		var data []byte
		data, _, err = handle.ReadPacketData()
		if err != nil && err != pcap.NextErrorTimeoutExpired || len(data) <= 0 {
			continue
		}
		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)
		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			continue
		}
		receiveTCP4 := packet.TransportLayer().(*layers.TCP)
		ack = receiveTCP4.Ack
		if ack != 0 && uint16(receiveTCP4.SrcPort) == dstPort && uint16(receiveTCP4.DstPort) == srcPort {
			seq = receiveTCP4.Seq
			window = receiveTCP4.Window
			log.Debugf("收到对端回复 %s", layersTCP2String(dstIP.String(), srcIP.String(), receiveTCP4))
			break
		}
	}
	if ack == 0 {
		err = fmt.Errorf("无法杀死连接 %s，拿不到 ack #_#", layersTCP2String(srcIP.String(), dstIP.String(), &tcp))
		return
	}

	tcp = layers.TCP{SrcPort: layers.TCPPort(srcPort), DstPort: layers.TCPPort(dstPort), Seq: ack, Ack: seq + 1, RST: true, SYN: true, Window: window}
	if err = tcp.SetNetworkLayerForChecksum(ip.(gopacket.NetworkLayer)); err != nil {
		return
	}
	buf = gopacket.NewSerializeBuffer()
	if err = gopacket.SerializeLayers(buf, options, &eth, ip, &tcp); err != nil {
		return
	}
	log.Debugf("发送数据包 [RST->remote] %s", layersTCP2String(srcIP.String(), dstIP.String(), &tcp))
	if err = handle.WritePacketData(buf.Bytes()); err != nil {
		return
	}

	srcIP, dstIP = dstIP, srcIP
	srcPort, dstPort = dstPort, srcPort
	tcp = layers.TCP{SrcPort: layers.TCPPort(srcPort), DstPort: layers.TCPPort(dstPort), Seq: seq, Ack: ack, RST: true, SYN: true, Window: window}
	if err = tcp.SetNetworkLayerForChecksum(ip.(gopacket.NetworkLayer)); err != nil {
		return
	}
	buf = gopacket.NewSerializeBuffer()
	if err = gopacket.SerializeLayers(buf, options, &eth, ip, &tcp); err != nil {
		return
	}
	log.Debugf("发送数据包 [RST->local]  %s", layersTCP2String(srcIP.String(), dstIP.String(), &tcp))
	if err = handle.WritePacketData(buf.Bytes()); err != nil {
		return
	}

	return
}

func layersTCP2String(srcIP, dstIP string, tcp *layers.TCP) string {
	s := ""
	if tcp.FIN {
		s += "FIN|"
	}
	if tcp.SYN {
		s += "SYN|"
	}
	if tcp.RST {
		s += "RST|"
	}
	if tcp.PSH {
		s += "PSH|"
	}
	if tcp.ACK {
		s += "ACK|"
	}
	if tcp.URG {
		s += "URG|"
	}
	if tcp.ECE {
		s += "ECE|"
	}
	if tcp.CWR {
		s += "CWR|"
	}
	if tcp.NS {
		s += "NS|"
	}
	return fmt.Sprintf("|src=%s:%d|dst=%s:%d|seq=%d|ack=%d|window=%d|checksum=%d|%s",
		srcIP, tcp.SrcPort, dstIP, tcp.DstPort, tcp.Seq, tcp.Ack, tcp.Window, tcp.Checksum, s,
	)
}
