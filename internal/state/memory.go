package state

import (
	"context"
	"sync"
)

type MemoryStore struct {
	mu          sync.RWMutex
	invocations map[string]*Invocation
}

func (m *MemoryStore) Close() error {
	m.Close()
	return nil
}

func (m *MemoryStore) MostRecentInvocation(ctx context.Context, collector string) (*Invocation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	i, found := m.invocations[collector]
	if !found {
		return nil, nil
	}

	return i, nil
}

func (m *MemoryStore) SaveInvocation(invocation *Invocation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	mr, found := m.invocations[invocation.CollectorName]

	if !found {
		m.invocations[invocation.CollectorName] = invocation
		return nil
	}

	// check if found infication is more recent
	if invocation.Time.After(mr.Time) {
		m.invocations[invocation.CollectorName] = invocation
	}

	return nil
}

func NewMemoryStoreFromGenericConfig(m map[string]any) (*MemoryStore, error) {
	s := &MemoryStore{
		invocations: make(map[string]*Invocation),
	}

	return s, nil
}
