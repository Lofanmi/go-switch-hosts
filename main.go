package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Lofanmi/go-switch-hosts/contracts"
	"github.com/Lofanmi/go-switch-hosts/internal/network"
	"github.com/Lofanmi/go-switch-hosts/internal/pcap"
	"github.com/Lofanmi/go-switch-hosts/internal/store"
	"github.com/google/wire"
)

type Application struct {
	HostsStore contracts.HostsStore
	Network    contracts.Network
}

var Sets = wire.NewSet(
	wire.Struct(new(Application), "*"),

	wire.Struct(new(store.ConfigLoader), "*"),
	wire.Bind(new(contracts.HostsConfigLoader), new(*store.ConfigLoader)),

	network.NewNetwork,
	pcap.NewHandleManager,
	store.NewHosts,
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	application, cleanup, err := NewApplication()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	gateway, iface, _ := application.Network.Route(net.ParseIP("10.43.21.88"))
	addr, err := application.Network.GatewayHardwareAddr(gateway, iface)
	fmt.Println(addr)

	c, _ := application.Network.GetTCPConnectionList()
	fmt.Println(c)
}
