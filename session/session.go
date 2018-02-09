package session

import (
	"github.com/satori/go.uuid"
	"log"
	"time"
)

const (
	keyExpires = "_expires"
)

type Session struct {
	// unique session id, used as redis key
	id     uuid.UUID
	values map[string]string
}

func New() (*Session, error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Printf("[session/session.go] UUID(version 4) generation failed\n")
		return nil, err
	}

	sess := &Session{
		id:     id,
		values: make(map[string]string),
	}
	sess.SetExpireAfter(time.Hour * 2)
	return sess, nil
}

func (sess *Session) ID() string {
	return sess.id.String()
}

// Get value
func (sess *Session) GetVal(key string) (string, bool) {
	val, ok := sess.values[key]
	return val, ok
}

// Set value
func (sess *Session) SetVal(key string, val string) {
	sess.values[key] = val
}

// Delete value
func (sess *Session) DelVal(key string) {
	delete(sess.values, key)
}

func (sess *Session) Expired() bool {
	val, ok := sess.GetVal(keyExpires)
	if ok == false {
		return true
	}
	expire, err := time.Parse(time.RFC1123, val)
	if err != nil || expire.Before(time.Now()) {
		return true
	}
	return false
}

func (sess *Session) SetExpire(t time.Time) {
	sess.values[keyExpires] = t.Format(time.RFC1123)
}

func (sess *Session) SetExpireAfter(d time.Duration) {
	sess.SetExpire(time.Now().Add(d))
}
