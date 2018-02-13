package user

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

type Repository interface {
	Add(user *User) (uint64, error)
	Update(user *User) error
	Remove(user *User) error
	Query(user *User) []User
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (repo *repository) Query(user *User) []User {
	users := []User{}
	if err := repo.db.Where(user).Find(&users).Error; err != nil {
		log.Printf("[user/repository.go] Querying user %s failed: %s\n", *user, err.Error())
		return nil
	}
	return users
}

func (repo *repository) Remove(user *User) error {
	if err := repo.db.Delete(user).Error; err != nil {
		log.Printf("[user/repository.go] Deleting user %v failed: %s\n", *user)
		return err
	}
	return nil
}

func (repo *repository) Add(user *User) (uint64, error) {
	if err := repo.db.Create(user).Error; err != nil {
		log.Printf("[user/repository.go] User %s already exists: %s\n", user.Name, err.Error())
		return 0, errors.New("User already exists")
	}
	return user.ID, nil
}

func (repo *repository) Update(user *User) error {
	if err := repo.db.Save(user).Error; err != nil {
		log.Printf("[user/repository.go] Updating user %s failed: %s\n", user.Name, err.Error())
		return err
	}
	return nil
}
