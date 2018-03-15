package user

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

var (
	ErrUserExists = errors.New("User already exists")
)

type Repository interface {
	Add(user *User) (*User, error)
	Update(user *User) error
	Patch(data map[string]interface{}) error
	Remove(user *User) error
	Query(user *User) ([]User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (repo *repository) Add(user *User) (*User, error) {
	if err := repo.db.Create(user).Error; err != nil {
		log.Debugf("[user/repository.go:Add] User %s already exists: %s\n", user.Name, err.Error())
		return nil, ErrUserExists
	}
	return user, nil
}

func (repo *repository) Update(user *User) error {
	if err := repo.db.Save(user).Error; err != nil {
		log.Debugf("[user/repository.go:Update] Updating user %s failed: %s\n", user.Name, err.Error())
		return err
	}
	return nil
}

func (repo *repository) Patch(data map[string]interface{}) error {
	if err := repo.db.Model(&User{}).Updates(data).Error; err != nil {
		log.Debugf("[user/repository.go:Patch] Patching user %v failed: %s\n", data, err.Error())
		return err
	}
	return nil
}

func (repo *repository) Remove(user *User) error {
	if err := repo.db.Delete(user).Error; err != nil {
		log.Debugf("[user/repository.go:Remove] Deleting user %v failed: %s\n", *user)
		return err
	}
	return nil
}

func (repo *repository) Query(user *User) ([]User, error) {
	users := []User{}
	if err := repo.db.Where(user).Find(&users).Error; err != nil {
		log.Debugf("[user/repository.go:Query] Querying user %s failed: %s\n", user, err.Error())
		return nil, err
	}
	log.Printf("data: %v\n", users)
	return users, nil
}
