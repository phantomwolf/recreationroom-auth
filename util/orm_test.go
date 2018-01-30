package util

import (
	"github.com/PhantomWolf/recreationroom-auth/model"
	"testing"
)

func TestDB(t *testing.T) {
	orm := ORM()
	users := []model.User{}
	orm.Find(&users)
	if len(users) == 0 {
		t.Fatal("Failed to find user")
	}
	t.Log(users)
}
