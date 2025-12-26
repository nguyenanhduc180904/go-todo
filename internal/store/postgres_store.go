package store

import (
	"database/sql"
	"go-todo/internal/model"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// InitSchema tạo bảng nếu chưa có (rất tiện khi deploy lần đầu)
func (s *PostgresStore) InitSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL
	);`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) Create(t model.Todo) error {
	query := `INSERT INTO todos (id, title, completed) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, t.ID, t.Title, t.Completed)
	return err
}

func (s *PostgresStore) GetAll() ([]model.Todo, error) {
	query := `SELECT id, title, completed FROM todos`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var t model.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func (s *PostgresStore) GetByID(id string) (model.Todo, error) {
	query := `SELECT id, title, completed FROM todos WHERE id = $1`
	var t model.Todo
	// QueryRow dùng cho truy vấn trả về 1 dòng
	err := s.db.QueryRow(query, id).Scan(&t.ID, &t.Title, &t.Completed)
	return t, err
}

func (s *PostgresStore) Update(id string, t model.Todo) error {
	query := `UPDATE todos SET title = $1, completed = $2 WHERE id = $3`
	_, err := s.db.Exec(query, t.Title, t.Completed, id)
	return err
}

func (s *PostgresStore) Delete(id string) error {
	query := `DELETE FROM todos WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}
