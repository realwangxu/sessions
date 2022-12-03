package sessions

type Session interface {
	Set(key, value interface{})
	Get(key interface{}) (interface{}, error)
	Remove(key interface{}) error
	SessionID() string
}
