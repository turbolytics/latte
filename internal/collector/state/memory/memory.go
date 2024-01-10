package memory

import (
	"github.com/turbolytics/collector/internal/collector/state"
	"sync"
)

type Store struct {
	mu sync.Mutex
	s  map[string][]*state.Invocation
}

func (m *Store) LastInvocation(collector string) (*state.Invocation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	is, found := m.s[collector]
	if !found {
		return nil, nil
	}

	return is[0], nil
}

func (m *Store) SaveInvocation(invocation *state.Invocation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.s[invocation.Collector] = append(m.s[invocation.Collector], invocation)
	return nil
}

func NewFromGenericConfig(m map[string]any) (*Store, error) {
	s := &Store{
		s: make(map[string][]*state.Invocation),
	}
	return s, nil
}
