package sessions

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type CookieManager struct {
	sync.RWMutex
	lifetime   time.Duration
	maxAge     int
	provider   Provider
	cookieName string
}

func NewCookieManager(providerName string, lifetime time.Duration, cookieName string, maxAge int) (*CookieManager, error) {
	provider, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", providerName)
	}
	return &CookieManager{
		lifetime:   lifetime,
		maxAge:     maxAge,
		provider:   provider,
		cookieName: cookieName,
	}, nil
}

func (m *CookieManager) Start(w http.ResponseWriter, r *http.Request) Session {
	m.Lock()
	defer m.Unlock()

	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		sid, _ := NewUUID()
		session, _ := m.provider.Init(sid)
		http.SetCookie(w, &http.Cookie{
			Name:     m.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   m.maxAge,
		})
		return session
	}
	sid, _ := url.QueryUnescape(cookie.Value)
	if session, err := m.provider.Read(sid); err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:     m.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   m.maxAge,
		})
		return session
	}
	sid, _ = NewUUID()
	session, _ := m.provider.Init(sid)
	http.SetCookie(w, &http.Cookie{
		Name:     m.cookieName,
		Value:    url.QueryEscape(sid),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   m.maxAge,
	})
	return session
}

func (m *CookieManager) Destory(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}

	m.Lock()
	defer m.Unlock()

	sid, _ := url.QueryUnescape(cookie.Value)
	m.provider.Destory(sid)
	http.SetCookie(w, &http.Cookie{
		Name:     m.cookieName,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Now(),
	})
}

func (m *CookieManager) GC(maxlifetime time.Duration) {
	m.Lock()
	defer m.Unlock()

	m.provider.GC(m.lifetime)
	time.AfterFunc(maxlifetime, func() { m.GC(maxlifetime) })
}
