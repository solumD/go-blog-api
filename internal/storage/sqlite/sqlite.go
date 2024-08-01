package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/solumD/go-blog-api/internal/types"
)

type Storage struct {
	db *sql.DB
}

// New создает новое sqlite хранилище
func New(path string) (*Storage, error) {
	const fnNew = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fnNew, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", fnNew, err)
	}

	return &Storage{db: db}, nil
}

// IsUserExist проверяет, есть ли в БД пользователь с указанным логином
func (s *Storage) IsUserExist(ctx context.Context, login string) (bool, error) {
	const fnIsUserExist = "storage.sqlite.IsUserExist"

	q := `SELECT COUNT(*) FROM users WHERE login = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, &login).Scan(&count); err != nil {
		return false, fmt.Errorf("%s: failed to check if user exists: %w", fnIsUserExist, err)
	}

	return count > 0, nil
}

// GetPassword получает пароль конкретного пользователя
func (s *Storage) GetPassword(ctx context.Context, login string) (string, error) {
	const fnGetPassword = "storage.sqlite.GetUserPassword"

	query := `SELECT password FROM users where login = ?`
	row := s.db.QueryRowContext(ctx, query, login)

	var password string

	err := row.Scan(&password)
	if err != nil {
		return "", fmt.Errorf("%s: failed to get %s's password: %w", fnGetPassword, login, err)
	}
	return password, nil
}

// SaveUser сохраняет пользователя и его зашифрованный пароль
func (s *Storage) SaveUser(ctx context.Context, login string, password string, date_registered string) (int64, error) {
	const fnSaveUser = "storage.sqlite.SaveUser"

	q := `
		INSERT INTO users(login, password, date_registered) VALUES(?, ?, ?)
	`
	data := []any{login, password, date_registered}

	res, err := s.db.ExecContext(ctx, q, data...)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to save user: %w", fnSaveUser, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert user's id: %w", fnSaveUser, err)
	}

	return id, nil
}

func (s *Storage) IsPostExist(ctx context.Context, id int) (bool, error) {
	const fnIsPostExist = "storage.sqlite.IsUserExist"

	q := `SELECT COUNT(*) FROM posts WHERE id = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, &id).Scan(&count); err != nil {
		return false, fmt.Errorf("%s: failed to check if post exists: %w", fnIsPostExist, err)
	}

	return count > 0, nil
}

// GetPostCreator получает id создателя поста
func (s *Storage) GetPostCreator(ctx context.Context, id int) (string, error) {
	const fnGetPostCreator = "storage.sqlite.GetPostCreator"

	q := `SELECT created_by FROM posts WHERE id = ?`

	var created_by string

	err := s.db.QueryRowContext(ctx, q, &id).Scan(&created_by)
	if err == sql.ErrNoRows {
		return "", err
	} else if err != nil {
		return "", fmt.Errorf("%s: failed to check if post exists: %w", fnGetPostCreator, err)
	}

	return created_by, nil
}

// GetPosts получает посты конкретного пользователя
func (s *Storage) GetPosts(ctx context.Context, created_by string) (*types.UsersPosts, error) {
	const fnGetPosts = "storage.sqlite.GetPosts"

	q := `
		SELECT posts.id, created_by, title, text, likes, date_created, date_updated FROM posts 
		WHERE created_by = ?
		ORDER BY posts.date_created desc`

	rows, err := s.db.QueryContext(ctx, q, created_by)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get %s's posts: %w", fnGetPosts, created_by, err)
	}
	defer rows.Close()

	posts := make([]types.Post, 0)
	for rows.Next() {
		var post types.Post
		if err := rows.Scan(&post.ID, &post.Created_by, &post.Title, &post.Text, &post.Likes, &post.Created_at, &post.Updated_at); err != nil {
			return nil, fmt.Errorf("%s: failed to scan %s's posts: %w", fnGetPosts, created_by, err)
		}
		posts = append(posts, post)
	}

	UserPosts := types.UsersPosts{Posts: posts}

	return &UserPosts, nil
}

// SavePost сохраняет пост пользователя
func (s *Storage) SavePost(ctx context.Context, created_by string, title string, text string, date_created string) (int64, error) {
	const fnSavePost = "storage.sqlite.SavePost"

	q := `
        INSERT INTO posts(created_by, title, text, date_created, date_updated) VALUES(?,?,?,?,?)`

	data := []any{created_by, title, text, date_created, date_created}

	res, err := s.db.ExecContext(ctx, q, data...)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to save post: %w", fnSavePost, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert post's id: %w", fnSavePost, err)
	}

	return id, nil
}

// UpdatePostTitle обновляет название поста
func (s *Storage) UpdatePostTitle(ctx context.Context, id int, title string, date_updated string) error {
	const fnUpdatePostTitle = "storage.sqlite.UpdatePostTitle"

	q := `UPDATE posts SET 
				title = ?,
				date_updated = ?
			  WHERE id = ?

					`
	data := []any{title, date_updated, id}

	_, err := s.db.ExecContext(ctx, q, data...)
	if err != nil {
		return fmt.Errorf("%s: failed to save post: %w", fnUpdatePostTitle, err)
	}

	return nil
}

// UpdatePostTest обновляет текст поста
func (s *Storage) UpdatePostText(ctx context.Context, id int, text string, date_updated string) error {
	const fnUpdatePostText = "storage.sqlite.UpdatePostText"

	q := `UPDATE posts SET 
				text = ?,
				date_updated = ?
			  WHERE id = ?

					`
	data := []any{text, date_updated, id}

	_, err := s.db.ExecContext(ctx, q, data...)
	if err != nil {
		return fmt.Errorf("%s: failed to save post: %w", fnUpdatePostText, err)
	}

	return nil
}

// RemovePost удаляет пост
func (s *Storage) RemovePost(ctx context.Context, id int) error {
	const fnRemovePost = "storage.sqlite.RemovePost"

	q := `DELETE FROM posts WHERE id = ?`

	_, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("%s: failed to delete post: %w", fnRemovePost, err)
	}

	return nil
}

func (s *Storage) IsPostLikedByUser(ctx context.Context, id int, liked_by string) (bool, error) {
	const fnIsPostLikedByUser = "storage.sqlite.IsPostLikedByUser"

	q := `SELECT COUNT(*) FROM reactions WHERE post_id = ? AND liked_by = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, &id, &liked_by).Scan(&count); err != nil {
		return false, fmt.Errorf("%s: failed to check if user exists: %w", fnIsPostLikedByUser, err)
	}

	return count > 0, nil
}

func (s *Storage) LikePost(ctx context.Context, id int, liked_by string) error {
	const fnLikePost = "storage.sqlite.LikePost"

	q := `UPDATE posts SET likes = likes + 1 WHERE id = ?`

	_, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("%s: failed to increase post's likes: %w", fnLikePost, err)
	}

	q = `INSERT INTO reactions (post_id, liked_by) VALUES(?, ?)`

	_, err = s.db.ExecContext(ctx, q, id, liked_by)
	if err != nil {
		return fmt.Errorf("%s: failed to save %s's like on post: %w", fnLikePost, liked_by, err)
	}

	return nil

}

func (s *Storage) UnlikePost(ctx context.Context, id int, liked_by string) error {
	const fnUnlikePost = "storage.sqlite.UnlikePost"

	q := `UPDATE posts SET likes = likes - 1 WHERE id = ?`

	_, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("%s: failed to decrease post's likes: %w", fnUnlikePost, err)
	}

	q = `DELETE FROM reactions WHERE post_id = ? AND liked_by = ?`

	_, err = s.db.ExecContext(ctx, q, id, liked_by)
	if err != nil {
		return fmt.Errorf("%s: failed to delete %s's like on post: %w", fnUnlikePost, liked_by, err)
	}

	return nil
}

// Init создает таблицы и индексы, если они еще не были созданы
func (s *Storage) Init(ctx context.Context) error {
	const fnInit = "storage.sqlite.Init"

	q := `
		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			login VARCHAR(50) NOT NULL, 
			password VARCHAR(143) NOT NULL,
			date_registered TIMESTAMP NOT NULL);

		CREATE INDEX IF NOT EXISTS idx_login ON users(login);
		
		CREATE TABLE IF NOT EXISTS posts(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_by VARCHAR(50) NOT NULL, 
			title VARCHAR(100) NOT NULL,
			text TEXT NOT NULL,
			likes INTEGER DEFAULT 0, 
			date_created TIMESTAMP NOT NULL, 
			date_updated TIMESTAMP NOT NULL,
			FOREIGN KEY(created_by) REFERENCES users(login) ON DELETE CASCADE);

		CREATE TABLE IF NOT EXISTS reactions(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			liked_by VARCHAR(50), 
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE);
		
	`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("%s: failed to init tables: %w", fnInit, err)
	}

	return nil
}
