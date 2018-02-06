package session

type Repository interface {
	Create(sess *Session) (*Session, error)
	Update(sess *Session) error
	Delete(sid string) error
	Query() (*Session, error)
}

type repository struct {
	storage Storage
}

func NewRepository(storage Storage) Repository {
	return &repository{storage: storage}
}

func (r *repository) Create(sess *Session) (*Session, error) {

}
