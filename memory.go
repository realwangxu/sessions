package sessions

import (
	"fmt"
	"sync"
	"time"
)

type Memory struct {
	sid              string
	lastAccessedTime time.Time
	data             map[interface{}]interface{}
}

func (m *Memory) Set(key, value interface{}) {
	m.data[key] = value
}

func (m *Memory) Get(key interface{}) (interface{}, error) {
	v, ok := m.data[key]
	if !ok {
		return nil, fmt.Errorf("memory is not found %v", key)
	}
	return v, nil
}

func (m *Memory) Remove(key interface{}) error {
	if _, ok := m.data[key]; !ok {
		return fmt.Errorf("memory remove key %v is not exists", key)
	}
	delete(m.data, key)
	return nil
}

func (m *Memory) SessionID() string {
	return m.sid
}

type MemoryStore struct {
	sync.RWMutex
	lifetime time.Duration
	sessions map[string]Session
}

func NewMemoryStore(lifetime time.Duration) *MemoryStore {
	return &MemoryStore{lifetime: lifetime, sessions: make(map[string]Session)}
}

func (m *MemoryStore) Init(sid string) (Session, error) {
	m.Lock()
	defer m.Unlock()

	elem := &Memory{sid: sid, lastAccessedTime: time.Now(), data: make(map[interface{}]interface{})}
	m.sessions[sid] = elem
	return elem, nil
}

func (m *MemoryStore) Read(sid string) (Session, error) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.sessions[sid]; !ok {
		return nil, fmt.Errorf("memory session id %v is not exists", sid)
	}
	s := m.sessions[sid]
	now := time.Now()
	if s.(*Memory).lastAccessedTime.Add(m.lifetime).Before(now) {
		delete(m.sessions, sid)
		return nil, fmt.Errorf("memory session id %v expired", sid)
	}
	s.(*Memory).lastAccessedTime = now
	return s, nil
}

func (m *MemoryStore) Destory(sid string) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.sessions[sid]; !ok {
		return fmt.Errorf("memory session id %v is not exists", sid)
	}
	delete(m.sessions, sid)
	return nil
}

func (m *MemoryStore) GC() {
	m.Lock()
	defer m.Unlock()

	var res []string

	current := time.Now()
	for i := range m.sessions {
		elem := m.sessions[i].(*Memory)
		if elem.lastAccessedTime.Add(m.lifetime).Before(current) {
			res = append(res, i)
		}
	}
	for i := range res {
		delete(m.sessions, res[i])
	}
}
