package events

import (
	"sync"
)

type PolicyUpdateListener chan PolicyEvent

type PolicyEventManager struct {
	mu        sync.Mutex
	listeners []PolicyUpdateListener
}

type PolicyEvent struct{}

func (m *PolicyEventManager) Subscribe() PolicyUpdateListener {
	m.mu.Lock()
	defer m.mu.Unlock()
	ch := make(PolicyUpdateListener, 1)
	m.listeners = append(m.listeners, ch)
	return ch
}

func (m *PolicyEventManager) Notify(event PolicyEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, ch := range m.listeners {
		select {
		case ch <- event:
		default:
		}
	}
}

func (m *PolicyEventManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.listeners = nil
}
