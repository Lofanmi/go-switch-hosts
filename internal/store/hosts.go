package store

import (
	"strings"

	"github.com/Lofanmi/go-switch-hosts/contracts"
	"github.com/Lofanmi/go-switch-hosts/internal/gotil"
	"github.com/elliotchance/orderedmap/v2"
)

type Hosts struct {
	loader contracts.HostsConfigLoader
	m      *orderedmap.OrderedMap[string, []string]
}

func NewHosts(loader contracts.HostsConfigLoader) contracts.HostsStore {
	return &Hosts{
		loader: loader,
		m:      orderedmap.NewOrderedMap[string, []string](),
	}
}

func (s *Hosts) Init() {
	path := s.loader.Path()
	s.loader.Load(path, s)
	s.loader.Print(s)
}

func (s *Hosts) Flush() {
	s.m = orderedmap.NewOrderedMap[string, []string]()
}

func (s *Hosts) Parse(content string) {
	for _, line := range strings.Split(strings.Trim(content, " \r\n\t"), "\n") {
		line = strings.ReplaceAll(strings.Trim(line, " \t"), "\t", " ")
		if line == "" || line[0] == '#' || line[0] == ';' {
			continue
		}
		pieces := strings.SplitN(line, " ", 2)
		if len(pieces) <= 1 || pieces[0] == "" {
			continue
		}
		IP, hosts := pieces[0], strings.Fields(pieces[1])
		for _, host := range hosts {
			s.Bind(IP, host)
		}
	}
}

func (s *Hosts) Query(IP string) (hosts []string) {
	hosts, _ = s.m.Get(IP)
	return
}

func (s *Hosts) Bind(IP, host string) {
	if IP == "" || host == "" {
		return
	}
	values, ok := s.m.Get(IP)
	if ok && gotil.InArray[string](host, values) {
		return
	}
	values = append(values, host)
	s.m.Set(IP, values)
}

func (s *Hosts) IPs() []string    { return s.m.Keys() }
func (s *Hosts) Forget(IP string) { s.m.Delete(IP) }
func (s *Hosts) Count() int       { return s.m.Len() }
