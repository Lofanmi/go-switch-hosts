package contracts

import (
	"io"

	"github.com/fsnotify/fsnotify"
)

const (
	ChangeTypeConfig ChangeType = iota
	ChangeTypeHosts
)

const (
	EntryTypeComment EntryType = iota
	EntryTypeEmpty
	EntryTypeIPHost
)

type (
	EntryType int

	IP   = string
	Host = string

	ChangeType  = int
	ChangeEvent struct {
		Type  ChangeType
		Event fsnotify.Event
	}

	Entry struct {
		Type    EntryType `json:"type"`
		Content string    `json:"content"`
		IP      IP        `json:"ip"`
		Host    Host      `json:"host"`
	}

	EntrySlice []Entry

	Parser interface {
		Parse(list *EntrySlice, content string)
		Bind(list *EntrySlice, IP, host string)
		Comment(list *EntrySlice, comment string)
		EmptyLine(list *EntrySlice)
	}

	HostsStore interface {
		Init()
		List() EntrySlice
		Save(list EntrySlice)
		Write(buf io.Writer) (err error)
	}

	HostsConfigLoader interface {
		Path() (path string)
		Load(path string)
		OnChange(func(event ChangeEvent))
	}
)

func (s Entry) String() string {
	switch s.Type {
	case EntryTypeIPHost:
		return s.Content
	case EntryTypeComment:
		return "# " + s.Content
	}
	return ""
}

func (s EntrySlice) Map() (m map[Host]IP) {
	if len(s) <= 0 {
		return
	}
	m = map[Host]IP{}
	for _, entry := range s {
		if entry.Type != EntryTypeIPHost {
			continue
		}
		// hosts 文件优先使用前面的定义，后面有相同主机名称但是 IP 不同也会被忽略。
		if _, exist := m[entry.Host]; !exist {
			m[entry.Host] = entry.IP
		}
	}
	return
}

func (s EntrySlice) Diff(newer EntrySlice) (change bool, older *Set) {
	m := s.Map()
	if len(m) <= 0 && len(newer) > 0 {
		change = true
		return
	}
	older = NewSet()
	for _, newEntry := range newer {
		if oldIP, exist := m[newEntry.Host]; exist && oldIP != newEntry.IP {
			change = true
			older.Add(oldIP)
		}
	}
	n := newer.Map()
	for host, ip := range m {
		if _, exist := n[host]; !exist {
			change = true
			older.Add(ip)
		}
	}
	return
}
