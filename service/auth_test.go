package service

import (
	"context"
	"testing"
)

func TestRegister(t *testing.T) {
	service := &auth{}
	ctx := context.Background()

	name := "nobody"
	password := "redhat"
	email := "nobody@example.com"
	uid, err := service.Register(ctx, name, password, email)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	t.Logf("User %s successfully registered, uid %d\n", name, uid)
}
