package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) SaveUser(
	ctx context.Context,
	email string,
	passHash []byte,
) (uid int64, err error) {

	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES (?, ?)")

	if err != nil {
		fmt.Errorf("%s, %w", op, err)
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, email, passHash)

	if err != nil {
		var sqliteErr sqlite3.Error

	}

}
