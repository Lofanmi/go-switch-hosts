package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Lofanmi/go-switch-hosts/contracts"
	"github.com/Lofanmi/go-switch-hosts/internal/network"
	"github.com/Lofanmi/go-switch-hosts/internal/pcap"
	"github.com/Lofanmi/go-switch-hosts/internal/store"
	"github.com/google/wire"
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
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
