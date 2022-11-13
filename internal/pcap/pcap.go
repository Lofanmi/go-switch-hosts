package pcap

import (
	"github.com/Lofanmi/go-switch-hosts/contracts"
	"github.com/google/gopacket/pcap"
)

type HandleManager struct {
	pcapHandles map[string]*pcap.Handle
}

func NewHandleManager() (manager contracts.PCAPHandleManager, fn func(), err error) {
	m := &HandleManager{pcapHandles: map[string]*pcap.Handle{}}
	fn = func() {
		for _, handle := range m.pcapHandles {
			handle.Close()
		}
	}
	manager = m
	return
}

func (s *HandleManager) GetHandle(device string) (handle *pcap.Handle, err error) {
	if v, ok := s.pcapHandles[device]; ok {
		handle = v
		return
	}
	if handle, err = pcap.OpenLive(device, 65536, true, pcap.BlockForever); err != nil {
		return
	}
	s.pcapHandles[device] = handle
	return
}
