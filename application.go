package main

import (
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Lofanmi/go-switch-hosts/contracts"
	"github.com/Lofanmi/go-switch-hosts/internal/network"
	"github.com/Lofanmi/go-switch-hosts/internal/pcap"
	"github.com/Lofanmi/go-switch-hosts/internal/store"
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"
)

var Sets = wire.NewSet(
	wire.Struct(new(Application), "*"),
	network.NewNetwork,
	pcap.NewHandleManager,
	store.NewConfigLoader,
	store.NewHostsStore,
	store.NewDefaultParser,
)

type Application struct {
	HostsStore contracts.HostsStore
	Network    contracts.Network
}

func (s Application) Run() {
	s.HostsStore.Init()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	for {
		switch <-ch {
		case syscall.SIGINT, syscall.SIGTERM:
			return
		case syscall.SIGUSR1:
			killLocalTest(s.Network, "22")
		}
	}
}

// 用于本地测试 (macOS)
func killLocalTest(network contracts.Network, port string) {
	cmd := exec.Command("sh", "-c", "netstat -na|grep 192.168|grep "+port+"|grep ESTABLISHED")
	data, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	str2uint16 := func(s string) uint16 {
		i, _ := strconv.ParseUint(s, 10, 64)
		return uint16(i)
	}
	piece2hp := func(s string) (host, port string) {
		s = strings.ReplaceAll(s, ".", ":")
		s = strings.Replace(s, ":", ".", 3)
		var e error
		host, port, e = net.SplitHostPort(s)
		if e != nil {
			log.Println(e)
		}
		return
	}
	tcp := new(contracts.TCPConnection)
	var ok bool
	for _, line := range strings.Split(strings.Trim(string(data), " \r\n\t"), "\n") {
		line = strings.ReplaceAll(strings.Trim(line, " \t"), "\t", " ")
		if strings.Contains(line, "."+port+" ") {
			pieces := strings.Fields(line)
			if len(pieces) <= 1 || pieces[0] == "" {
				return
			}
			srcIP, srcPort := piece2hp(pieces[3])
			dstIP, dstPort := piece2hp(pieces[4])
			tcp.SrcIP = net.ParseIP(srcIP)
			tcp.SrcPort = str2uint16(srcPort)
			tcp.DstIP = net.ParseIP(dstIP)
			tcp.DstPort = str2uint16(dstPort)
			log.Debugf("端口[%s]存在连接 %s", port, tcp.String())
			ok = true
			break
		}
	}
	if !ok {
		log.Debugf("未发现该端口[%s]的连接 ，无法测试", port)
		return
	}
	err = network.KillTCPConnection(tcp)
	if err != nil {
		log.Debugf("%+v", err)
	}
}
