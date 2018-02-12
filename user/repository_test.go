package user

import (
	"github.com/PhantomWolf/recreationroom-auth/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"testing"
)

func TestRepositoryQuery(t *testing.T) {
	config.Load()
	db, err := gorm.Open(config.DatabaseBackend(), config.DSN())
	if err != nil {
		t.Fatalf("Database connection failed: %s\n", err.Error())
	}
	defer func() {
		db.Unscoped().Where("name = ?", "foo").Delete(&User{})
		db.Close()
	}()

	user := &User{Name: "foo", Password: "bar", Email: "foobar@example.com"}
	db.Create(user)
	t.Logf("Data prepared: %v\n", *user)

	repo := NewRepository(db)
	t.Log("[Test 1] Querying by user name:")
	users := repo.Query(&User{Name: "foo"})
	if len(users) != 1 {
		t.Fatal("No user found")
	}

	t.Log("[Test 2] Query by user email:")
	users = repo.Query(&User{Email: "foobar@example.com"})
	for _, u := range users {
		t.Logf("\tID: %d, Name: %s, Password: %s, Email: %s\n", u.ID, u.Name, u.Password, u.Email)
	}
	if len(users) == 0 {
		t.Fatal("No user found")
	}

	t.Log("[Test 3] Querying non-existing user:")
	users = repo.Query(&User{Name: "nobody"})
	if len(users) != 0 {
		t.Fatalf("User found unexpectedly: %v\n", users)
	}
}

func TestRepositoryAdd(t *testing.T) {
	config.Load()
	db, err := gorm.Open(config.DatabaseBackend(), config.DSN())
	if err != nil {
		t.Fatalf("Database connection failed: %s\n", err.Error())
	}
	defer func() {
		db.Unscoped().Where("name = ?", "foo").Delete(&User{})
		db.Close()
	}()

	user, _ := New("foo", "bar", "foo@example.com")
	repo := NewRepository(db)
	t.Log("[Test 1] Adding user")
	uid, err := repo.Add(user)
	if err != nil {
		t.Fatalf("Adding user failed: %s\n", err.Error())
	}
	// See if the user has been added
	out := &User{}
	db.Where(&User{Name: "foo"}).First(out)
	if out.ID != uid {
		t.Fatalf("Incorrect user id retrived: %d\n", uid)
	} else {
		t.Logf("User succesfully added, uid: %d\n", out.ID)
	}
}

func TestRepositoryUpdate(t *testing.T) {
	config.Load()
	db, err := gorm.Open(config.DatabaseBackend(), config.DSN())
	if err != nil {
		t.Fatalf("Database connection failed: %s\n", err.Error())
	}
	defer func() {
		db.Unscoped().Where("name = ?", "foo").Delete(&User{})
		db.Close()
	}()

}
