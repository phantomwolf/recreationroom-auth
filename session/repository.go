package session

import (
	"errors"
	"github.com/satori/go.uuid"
	"log"
)

type Repository interface {
	Add(sess *Session) error
	Update(sess *Session) error
	Remove(sid string) error
	Find(sid string) (*Session, error)
}

type repository struct {
	storage Storage
}

func NewRepository(storage Storage) Repository {
	return &repository{storage: storage}
}

func (repo *repository) Find(sid string) (*Session, error) {
	// Make sure session id is a valid UUID
	guid, err := uuid.FromString(sid)
	if err != nil {
		log.Printf("[session/repository.go] Invalid session id %s: %s\n", sid, err.Error())
		return nil, errors.New("Invalid session id")
	}

	data, err := repo.storage.Load(sid)
	if err != nil {
		return nil, err
	}
	sess := &Session{
		id:     guid,
		values: data,
	}
	return sess, nil
}

func (repo *repository) Remove(sid string) error {
	err := repo.storage.Delete(sid)
	return err
}

func (repo *repository) Update(sess *Session) error {
	id := sess.id.String()
	err := repo.storage.Save(id, sess.values)
	if err != nil {
		log.Printf("[session/repository.go] Session %s saving failure\n", id)
		return err
	}
	return nil
}

func (repo *repository) Add(sess *Session) error {
	id := sess.id.String()
	if repo.storage.Exists(id) {
		log.Printf("[session/repository.go] Session %s already exists\n", id)
		return errors.New("Session already exists")
	}

	err := repo.storage.Save(id, sess.values)
	if err != nil {
		log.Printf("[session/repository.go] Session %s saving failure: %s\n", id, err.Error())
		return err
	}
	return nil
}
