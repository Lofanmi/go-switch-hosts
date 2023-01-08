package contracts

type Set struct {
	m map[string]struct{}
}

func NewSet() *Set {
	return &Set{m: map[string]struct{}{}}
}

func (s *Set) Len() int {
	return len(s.m)
}

func (s *Set) Add(item string) {
	s.m[item] = struct{}{}
}

func (s *Set) Exist(item string) bool {
	_, ok := s.m[item]
	return ok
}

func (s *Set) Range(fn func(item string) bool) {
	for e := range s.m {
		if !fn(e) {
			break
		}
	}
}
