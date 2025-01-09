package postgres

import (
	"GoNews/pkg/storage"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // Драйвер PostgreSQL
)

// Store - структура для работы с PostgreSQL.
type Store struct {
	db *sql.DB
}

// New создает новый экземпляр Store.
func New(connectionString string) (*Store, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Store{db: db}, nil
}

// Posts возвращает все публикации из базы данных.
func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query("SELECT p.id, a.name, p.title, p.content, p.created_at FROM posts p JOIN authors a ON p.author_id = a.id")
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		if err := rows.Scan(&p.ID, &p.AuthorName, &p.Title, &p.Content, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// AddPost добавляет новую публикацию в базу данных.
func (s *Store) AddPost(p storage.Post) error {
	//  Получаем ID автора или создаем нового, если его нет
	var authorID int
	err := s.db.QueryRow("SELECT id FROM authors WHERE name = $1", p.AuthorName).Scan(&authorID)
	if err == sql.ErrNoRows {
		err = s.db.QueryRow("INSERT INTO authors (name) VALUES ($1) RETURNING id", p.AuthorName).Scan(&authorID)
		if err != nil {
			return fmt.Errorf("failed to create author: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to get author ID: %w", err)
	}

	if p.CreatedAt == 0 {
		p.CreatedAt = time.Now().Unix()
	}

	_, err = s.db.Exec("INSERT INTO posts (author_id, title, content, created_at) VALUES ($1, $2, $3, $4)", authorID, p.Title, p.Content, p.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to add post: %w", err)
	}
	return nil
}

// UpdatePost обновляет публикацию в базе данных.
func (s *Store) UpdatePost(p storage.Post) error {
	_, err := s.db.Exec("UPDATE posts SET title = $1, content = $2 WHERE id = $3", p.Title, p.Content, p.ID)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

// DeletePost удаляет публикацию из базы данных.
func (s *Store) DeletePost(p storage.Post) error {
	_, err := s.db.Exec("DELETE FROM posts WHERE id = $1", p.ID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}
