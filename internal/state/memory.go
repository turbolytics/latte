package state

import (
	"context"
	"github.com/turbolytics/latte/internal/invoker"
	"go.uber.org/zap"
	"sync"
)

type MemoryStore struct {
	mu          sync.RWMutex
	invocations map[string]*invoker.Invocation

	logger *zap.Logger
}

func (m *MemoryStore) Close() error {
	m.Close()
	return nil
}

func (m *MemoryStore) MostRecentInvocation(ctx context.Context, collector string) (*invoker.Invocation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	i, found := m.invocations[collector]
	if !found {
		return nil, nil
	}

	return i, nil
}

func (m *MemoryStore) SaveInvocation(invocation *invoker.Invocation) error {
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

type MemoryStoreOption func(store *MemoryStore)

func MemoryStoreWithLogger(l *zap.Logger) func(store *MemoryStore) {
	return func(ms *MemoryStore) {
		ms.logger = l
	}
}

func NewMemoryStoreFromGenericConfig(m map[string]any, opts ...MemoryStoreOption) (*MemoryStore, error) {
	s := &MemoryStore{
		invocations: make(map[string]*invoker.Invocation),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}
