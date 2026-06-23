package users

import (
	"context"
	"database/sql"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

const createTableQuery = `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL
	)`

func (r *Repo) CreateTable(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, createTableQuery)
	return err
}

func (r *Repo) Insert(ctx context.Context, name string, age int) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (name, age) VALUES (?, ?)", name, age)
	return err
}

type User struct {
	Id   int
	Name string
	Age  int
}

func (r *Repo) List(ctx context.Context) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, age FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Age); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *Repo) Get(ctx context.Context, id int) (User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, age FROM users WHERE id = ?", id)

	var user User
	err := row.Scan(&user.Id, &user.Name, &user.Age)
	return user, err
}

func (r *Repo) Tran(ctx context.Context) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// selalu abis create tx
	defer tx.Rollback()

	// tx.Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
	// tx.Rollback()
	// either commit atau rollback harus ada di akhir
}