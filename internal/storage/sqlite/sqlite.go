package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/solumD/go-blog-api/internal/types"
)

type Storage struct {
	db *sql.DB
}

const (
	fnNew         = "storage.sqlite.New"
	fnSaveUser    = "storage.sqlite.SaveUser"
	fnIsUserExist = "storage.sqlite.IsUserExist"
	fnSavePost    = "storage.sqlite.SavePost"
	fnGetPosts    = "storage.sqlite.GetPosts"
	fnDeletePost  = "storage.sqlite.DeletePost"
	fnInit        = "storage.sqlite.Init"
)

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fnNew, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", fnNew, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) IsUserExist(login string) (bool, error) {
	q := `SELECT COUNT(*) FROM users where login = ?`

	var count int

	if err := s.db.QueryRow(q, &login).Scan(&count); err != nil {
		return false, fmt.Errorf("%s: failed to check if user exists: %w", fnIsUserExist, err)
	}

	return count > 0, nil
}

func (s *Storage) SaveUser(login string, password string) (int64, error) {
	q := `
		INSERT INTO users(login, password) VALUES(?, ?)
	`
	data := []any{login, password}

	res, err := s.db.Exec(q, data...)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to save user: %w", fnSaveUser, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert user's id: %w", fnSaveUser, err)
	}

	return id, nil
}

func (s *Storage) SavePost(created_by string, title string, text string, date_created string) (int64, error) {
	q := `
        INSERT INTO posts(created_by, title, text, date_created, date_updated) VALUES(?,?,?,?,?)
    `
	data := []any{created_by, title, text, date_created, date_created}

	res, err := s.db.Exec(q, data...)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to save post: %w", fnSavePost, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert post's id: %w", fnSavePost, err)
	}

	return id, nil
}

func (s *Storage) GetPosts(created_by string) (*types.UsersPosts, error) {
	query := `
		SELECT posts.id, created_by, title, text, date_created, date_updated FROM posts
		JOIN users
			ON users.login = posts.created_by
			WHERE login = ?
		ORDER BY posts.id desc`

	rows, err := s.db.Query(query, created_by)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get %s's posts: %w", fnGetPosts, created_by, err)
	}
	defer rows.Close()

	posts := make([]types.Post, 0)
	for rows.Next() {
		var post types.Post
		if err := rows.Scan(&post.ID, &post.Created_by, &post.Title, &post.Text, &post.Created_by, &post.Created_at); err != nil {
			return nil, fmt.Errorf("%s: failed to scan %s's posts: %w", fnGetPosts, created_by, err)
		}
		posts = append(posts, post)
	}

	UserPosts := types.UsersPosts{}
	UserPosts.Posts = posts

	return &UserPosts, nil
}

func (s *Storage) DeletePost(id int) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%s: failed to delete post: %w", fnDeletePost, err)
	}
	return nil
}

func (s *Storage) Init() error {
	q := `
		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			login VARCHAR(50) NOT NULL, 
			password VARCHAR(143) NOT NULL);

		CREATE INDEX IF NOT EXISTS idx_login ON users(login);
		
		CREATE TABLE IF NOT EXISTS posts(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_by VARCHAR(50) NOT NULL, 
			title VARCHAR(100) NOT NULL,
			text TEXT NOT NULL, 
			date_created TIMESTAMP NOT NULL, 
			date_updated TIMESTAMP NOT NULL,
			FOREIGN KEY(created_by) REFERENCES users(login) ON DELETE CASCADE);
	`
	_, err := s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("%s: failed to init tables: %w", fnInit, err)
	}

	return nil
}
