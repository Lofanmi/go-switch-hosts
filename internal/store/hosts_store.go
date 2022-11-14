package store

import (
	"os"
	"sync"

	"github.com/Lofanmi/go-switch-hosts/contracts"
	"github.com/Lofanmi/go-switch-hosts/internal/gotil"
)

type HostsStore struct {
	loader contracts.HostsConfigLoader
	parser contracts.Parser
	list   contracts.EntrySlice
	mu     *sync.RWMutex
}

func NewHostsStore(loader contracts.HostsConfigLoader, parser contracts.Parser) contracts.HostsStore {
	return &HostsStore{
		loader: loader,
		parser: parser,
		mu:     new(sync.RWMutex),
	}
}

func (s *HostsStore) Init() {
	path := s.loader.Path()
	s.loader.Load(path, s.parser)
	s.Save(s.parser.List())
	if gotil.IsDevelopment() {
		_ = s.loader.Print(s, os.Stdout)
	}
}

func (s *HostsStore) Save(list contracts.EntrySlice) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.list = list
}

func (s *HostsStore) List() contracts.EntrySlice {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n := len(s.list)
	result := make(contracts.EntrySlice, n)
	for i := 0; i < n; i++ {
		result[i] = s.list[i]
	}
	return result
}

func (s *HostsStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.list)
}

func (s *HostsStore) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.list = nil
}
