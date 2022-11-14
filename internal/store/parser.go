package store

import (
	"fmt"
	"net"
	"strings"

	"github.com/Lofanmi/go-switch-hosts/contracts"
)

type DefaultParser struct {
	list contracts.EntrySlice
}

func NewDefaultParser() contracts.Parser {
	return &DefaultParser{}
}

func (s *DefaultParser) Parse(content string) {
	for _, line := range strings.Split(strings.Trim(content, " \r\n\t"), "\n") {
		if line = strings.ReplaceAll(strings.Trim(line, " \t"), "\t", " "); line == "" {
			s.EmptyLine()
			continue
		}
		if line[0] == '#' || line[0] == ';' {
			s.Comment(line[1:])
			continue
		}
		pieces := strings.SplitN(line, " ", 2)
		if len(pieces) <= 1 || pieces[0] == "" {
			continue
		}
		IP, hosts := pieces[0], strings.Fields(pieces[1])
		if net.ParseIP(IP) == nil {
			s.Comment(fmt.Sprintf("[PARSE ERROR] %s", line))
			continue
		}
		for _, host := range hosts {
			s.Bind(IP, host)
		}
	}
}

func (s *DefaultParser) Bind(IP, host string) {
	IP = strings.TrimSpace(IP)
	host = strings.TrimSpace(host)
	if IP == "" || host == "" {
		return
	}
	s.list = append(s.list, contracts.Entry{
		Type:    contracts.EntryTypeIPHost,
		Content: fmt.Sprintf("%s %s", IP, host),
		IP:      IP,
		Host:    host,
	})
}

func (s *DefaultParser) Comment(comment string) {
	if comment = strings.TrimSpace(comment); comment == "" {
		return
	}
	s.list = append(s.list, contracts.Entry{Type: contracts.EntryTypeComment, Content: comment})
}

func (s *DefaultParser) EmptyLine() {
	s.list = append(s.list, contracts.Entry{Type: contracts.EntryTypeEmpty})
}

func (s *DefaultParser) List() (list contracts.EntrySlice) {
	return s.list
}
