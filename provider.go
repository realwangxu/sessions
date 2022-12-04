package sessions

import (
	"net/http"
	"time"
)

type Provider interface {
	Init(sid string) (Session, error)
	Read(sid string) (Session, error)
	Destory(sid string) error
	GC()
}

var (
	providers = make(map[string]Provider)
	global    = &CookieManager{}
)

func Register(name string, provider Provider) {
	if provider == nil {
		panic("provider register: provider is nil")
	}
	if _, ok := providers[name]; ok {
		panic("provider register: provider already exists")
	}
	providers[name] = provider
}

func WithBackground(manager *CookieManager) {
	global = manager
}

func Start(w http.ResponseWriter, r *http.Request) Session {
	return global.Start(w, r)
}

func Destory(w http.ResponseWriter, r *http.Request) {
	global.Destory(w, r)
}

func GC(lifetime time.Duration) {
	global.GC(lifetime)
}
