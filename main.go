package main

import (
	"fmt"
	"net"

	"github.com/Lofanmi/go-switch-hosts/contracts"
	"github.com/Lofanmi/go-switch-hosts/internal/gotil"
	"github.com/Lofanmi/go-switch-hosts/internal/network"
	"github.com/Lofanmi/go-switch-hosts/internal/pcap"
	"github.com/Lofanmi/go-switch-hosts/internal/store"
	"github.com/google/wire"
	log "github.com/sirupsen/logrus"
)

type Application struct {
	HostsStore contracts.HostsStore
	Network    contracts.Network
}

var Sets = wire.NewSet(
	wire.Struct(new(Application), "*"),

	network.NewNetwork,
	pcap.NewHandleManager,

	wire.Struct(new(store.ConfigLoader), "*"),
	wire.Bind(new(contracts.HostsConfigLoader), new(*store.ConfigLoader)),

	store.NewHostsStore,
	store.NewDefaultParser,
)

func main() {
	initLogger()

	application, cleanup, err := NewApplication()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	application.HostsStore.Init()

	gateway, iface, _ := application.Network.Route(net.ParseIP("10.43.21.88"))
	addr, err := application.Network.GatewayHardwareAddr(gateway, iface)
	fmt.Println(addr)

	c, _ := application.Network.GetTCPConnectionList()
	fmt.Println(c)
}

const defaultLogLevel = "debug"

func initLogger() {
	if level, err := log.ParseLevel(gotil.Env(gotil.EnvGoSwitchHostsLogLevel, defaultLogLevel)); err != nil {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(level)
	}
}
