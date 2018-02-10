package user

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Repository interface {
	Add(user *User) error
	Update(user *User) error
	Remove(id int64) error
	Query(spec *User) []User
}

type repository struct {
	db *gorm.DB
}

func (repo *repository) Query(spec *User) []User {
    users := []User
	repo.db.Where(spec).Find(&users)
}

func (repo *repository) Remove(id int64) error {
	if err := repo.db.Where("id = ?", id).Delete(&User{}).Error; err != nil {
		log.Printf("[user/repository.go] Deleting user %s failed: %s\n", id)
		return err
	}
	return nil
}

func (repo *repository) Add(user *User) error {
	if err := repo.db.Create(user).Error; err != nil {
		log.Printf("[user/repository.go] User %s already exists: %s\n", user.Name, err.Error())
		return errors.New("User already exists")
	}
	return nil
}

func (repo *repository) Update(user *User) error {
	if err := repo.db.Save(user).Error; err != nil {
		log.Printf("[user/repository.go] Updating user %s failed: %s\n", user.Name, err.Error())
		return err
	}
	return nil
}
