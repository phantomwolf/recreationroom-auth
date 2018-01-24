package dao

import (
	"context"
	"github.com/PhantomWolf/recreationroom-auth/model/entity"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestUserDao(t *testing.T) {
	ctx := context.Background()
	// Create UserDao
	userDao, err := NewUser("mysql")
	if err != nil {
		t.Fatal("Failed to create UserDao: %v\n", err)
	}
	// Create user
	user := &entity.User{Name: "testuser", Password: "testpassword", Email: "testuser@example.com"}
	id, err := userDao.Create(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v\n", err)
	}
	// Update user
	user.Name = "rookie"
	user.Password = "rookiepassword"
	user.Email = "rookie@example.com"
	err = userDao.Update(ctx, user)
	if err != nil {
		t.Fatalf("Failed to update user: %v\n", err)
	}
	// Read user
	data, err := userDao.Read(ctx, id)
	if err != nil {
		t.Fatalf("Failed to read user: %v\n", err)
	}
	if data.Name != user.Name ||
		data.Password != user.Password ||
		data.Email != user.Email {
		t.Fatalf("Data read doesn't match with input: %v <=> %v\n", data, user)
	}
	// Delete User
	err = userDao.Delete(ctx, id)
	if err != nil {
		t.Fatalf("Failed to delete user: %v\n")
	}
}
