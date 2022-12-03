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

type MemoryManager struct {
	sync.RWMutex
	lifetime time.Duration
	sessions map[string]Session
}

func (m *MemoryManager) Expire(lifetime time.Duration) {
	m.lifetime = lifetime
}

func (m *MemoryManager) Init(sid string) (Session, error) {
	m.Lock()
	defer m.Unlock()

	elem := &Memory{sid: sid, lastAccessedTime: time.Now(), data: make(map[interface{}]interface{})}
	m.sessions[sid] = elem
	return elem, nil
}

func (m *MemoryManager) Read(sid string) (Session, error) {
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

func (m *MemoryManager) Destory(sid string) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.sessions[sid]; !ok {
		return fmt.Errorf("memory session id %v is not exists", sid)
	}
	delete(m.sessions, sid)
	return nil
}

func (m *MemoryManager) GC() {
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
