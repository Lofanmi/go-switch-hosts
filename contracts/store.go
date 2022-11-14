package contracts

import (
	"io"
)

type EntryType int

const (
	EntryTypeComment EntryType = iota
	EntryTypeEmpty
	EntryTypeIPHost
)

type Entry struct {
	Type    EntryType `json:"type"`
	Content string    `json:"content"`
	IP      string    `json:"ip"`
	Host    string    `json:"host"`
}

type EntrySlice []Entry

func (s Entry) String() string {
	switch s.Type {
	case EntryTypeIPHost:
		return s.Content
	case EntryTypeComment:
		return "# " + s.Content
	}
	return ""
}

type Parser interface {
	Parse(content string)
	Bind(IP, host string)
	Comment(comment string)
	EmptyLine()
	List() (list EntrySlice)
}

type HostsStore interface {
	Init()
	List() EntrySlice
	Save(list EntrySlice)
	Count() int
	Flush()
}

type HostsConfigLoader interface {
	Path() (path string)
	Load(path string, parser Parser)
	Print(hosts HostsStore, buf io.Writer) (err error)
}
