package user

import (
	"regexp"

	_ "github.com/jinzhu/gorm"
)

type User struct {
	ID       uint64 `gorm:"primary_key"`
	Name     string `sql:"type:varchar(50);unique_index"`
	Password string `sql:"type:varchar(50)"`
	Email    string `sql:"type:varchar(255)"`
	Nickname string `sql:"type:"`
}

func (u *User) SetName(name string) error {
	matched, err := regexp.MatchString("^\\w+$", name)
	if err != nil || matched == false {
		return &ErrInvalidName{"Username can only contain alphabets and numbers"}
	}
	u.Name = name
	return nil
}

func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return &ErrInvalidPassword{"Empty password"}
	}
	for _, c := range password {
		if c < 32 || c > 126 {
			return &ErrInvalidPassword{"Password should contain only ASCII characters"}
		}
	}
	u.Password = password
	return nil
}

func (u *User) SetEmail(email string) error {
	pattern := "^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$"
	if matched, err := regexp.MatchString(pattern, email); err != nil {
		panic("Regular expression error")
	} else if matched == false {
		return &ErrInvalidEmail{"Invalid email"}
	}
	u.Email = email
	return nil
}

func (u *User) SetNickname(nickname string) error {
	if len(nickname) == 0 {
		return &ErrInvalidNickname{"Nickname can't be empty"}
	}
	u.Nickname = nickname
	return nil
}

func New(name string, password string, email string) (*User, error) {
	p := new(User)
	if err := p.SetName(name); err != nil {
		return nil, err
	}
	if err := p.SetPassword(password); err != nil {
		return nil, err
	}
	if err := p.SetEmail(email); err != nil {
		return nil, err
	}
	return p, nil
}

type ErrInvalidName struct {
	s string
}

func (e *ErrInvalidName) Error() string {
	return e.s
}

type ErrInvalidPassword struct {
	s string
}

func (e *ErrInvalidPassword) Error() string {
	return e.s
}

type ErrInvalidEmail struct {
	s string
}

func (e *ErrInvalidEmail) Error() string {
	return e.s
}

type ErrInvalidNickname struct {
	s string
}

func (e *ErrInvalidNickname) Error() string {
	return e.s
}
