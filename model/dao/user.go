package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/PhantomWolf/recreationroom-auth/model/entity"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type User interface {
	Create(ctx context.Context, user *entity.User) (int64, error)
	Read(ctx context.Context, id int64) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	//Delete(id string) error
}

type mysql struct {
}

func (userDao *mysql) Delete(ctx context.Context, id int64) error {
	db := ctx.Value("db").(*sql.DB)
	query := "DELETE FROM user WHERE id = ?"
	res, err := db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("[model.dao.user] Deleting user %d failed", id)
	}
	return err
}

func (userDao *mysql) Update(ctx context.Context, user *entity.User) error {
	db := ctx.Value("db").(*sql.DB)
	query := "UPDATE user SET name = ?, password = ?, email = ? WHERE id = ?"
	res, err := db.ExecContext(ctx, query, user.Name, user.Password, user.Email, user.Id)
	if err != nil {
		log.Printf("[model.dao.user] Updating user %d failed\n", user.Id)
	}
	return err
}

// Return id of created user and error
func (userDao *mysql) Create(ctx context.Context, user *entity.User) (int64, error) {
	db := ctx.Value("db").(*sql.DB)
	query := "INSERT INTO user(id, name, password, email) VALUES(NULL, ?, ?, ?)"
	res, err := db.ExecContext(ctx, query, user.Name, user.Password, user.Email)
	if err != nil {
		log.Printf("[model.dao.user] Insertion failed: %v\n", *user)
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("[model.dao.user] Failed to get LastInsertId\n")
		return -1, err
	}
	return id, nil
}

func (userDao *mysql) Read(ctx context.Context, id int64) (*entity.User, error) {
	db := ctx.Value("db").(*sql.DB)
	rows, err := db.QueryContext(ctx, "SELECT name, password, email FROM user WHERE id = ?", id)
	defer rows.Close()
	if err != nil {
		log.Printf("[model.dao.user] SQL query failed: %v\n", err)
		return nil, err
	}
	if !rows.Next() {
		log.Printf("[model.dao.user] No such user %s\n", id)
		return nil, errors.New("No such user " + string(id))
	}
	// Extract data
	user := &entity.User{Id: id}
	err = rows.Scan(&user.Name, &user.Password, &user.Email)
	if err != nil {
		log.Println("[model.dao.user] Failed to extract data from query result")
		return nil, err
	}
	return user, nil
}

func NewUser(database string) (User, error) {
	switch database {
	case "mysql":
		return &mysql{}, nil
	default:
		return nil, errors.New("[model.dao.user] Unsupported database: " + database)
	}
}
