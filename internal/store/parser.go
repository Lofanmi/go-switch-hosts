package store

import (
	"fmt"
	"net"
	"strings"

	"github.com/Lofanmi/go-switch-hosts/contracts"
)

type DefaultParser struct{}

func NewDefaultParser() contracts.Parser {
	return &DefaultParser{}
}

func (s *DefaultParser) Parse(list *contracts.EntrySlice, content string) {
	for _, line := range strings.Split(strings.Trim(content, " \r\n\t"), "\n") {
		if line = strings.ReplaceAll(strings.Trim(line, " \t"), "\t", " "); line == "" {
			s.EmptyLine(list)
			continue
		}
		if line[0] == '#' || line[0] == ';' {
			s.Comment(list, line[1:])
			continue
		}
		pieces := strings.SplitN(line, " ", 2)
		if len(pieces) <= 1 || pieces[0] == "" {
			continue
		}
		IP, hosts := pieces[0], strings.Fields(pieces[1])
		if net.ParseIP(IP) == nil {
			s.Comment(list, fmt.Sprintf("[PARSE ERROR] %s", line))
			continue
		}
		for _, host := range hosts {
			host = strings.TrimRight(host, ";") // 移除末尾的 ; 符号
			s.Bind(list, IP, host)
		}
	}
	return
}

func (s *DefaultParser) Bind(list *contracts.EntrySlice, IP, host string) {
	IP = strings.TrimSpace(IP)
	host = strings.TrimSpace(host)
	if IP == "" || host == "" {
		return
	}
	*list = append(*list, contracts.Entry{
		Type:    contracts.EntryTypeIPHost,
		Content: fmt.Sprintf("%s %s", IP, host),
		IP:      IP,
		Host:    host,
	})
}

func (s *DefaultParser) Comment(list *contracts.EntrySlice, comment string) {
	comment = strings.TrimSpace(comment)
	*list = append(*list, contracts.Entry{Type: contracts.EntryTypeComment, Content: comment})
}

func (s *DefaultParser) EmptyLine(list *contracts.EntrySlice) {
	*list = append(*list, contracts.Entry{Type: contracts.EntryTypeEmpty})
}
