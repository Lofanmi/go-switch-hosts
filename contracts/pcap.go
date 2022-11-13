package contracts

import (
	"github.com/google/gopacket/pcap"
)

type PCAPHandleManager interface {
	GetHandle(device string) (handle *pcap.Handle, err error)
}
