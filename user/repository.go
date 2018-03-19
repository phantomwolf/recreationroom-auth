package user

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
		return nil, ErrUserAlreadyExists
	}
	return user, nil
}

func (repo *repository) Update(user *User) error {
	if err := repo.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (repo *repository) Patch(data map[string]interface{}) error {
	if err := repo.db.Model(&User{}).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (repo *repository) Remove(user *User) error {
	if err := repo.db.Delete(user).Error; err != nil {
		return err
	}
	return nil
}

func (repo *repository) Query(user *User) ([]User, error) {
	users := []User{}
	if err := repo.db.Where(user).Find(&users).Error; err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	return users, nil
}
