package models

import "github.com/solumD/go-blog-api/internal/types"

// Это пакет сделан для генерации swagger-документации

// auth

type RegisterSuccess struct {
	Status string `json:"status"`
	ID     int64  `json:"id"`
}

type RegisterError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type LoginSuccess struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

type LoginError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// post

type SaveSuccess struct {
	Status string `json:"status"`
	ID     int64  `json:"id"`
}

type SaveError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type DeleteSuccess struct {
	Status string `json:"status"`
}

type DeleteError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type UpdateSuccess struct {
	Status string `json:"status"`
}

type UpdateError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type LikeSuccess struct {
	Status string `json:"status"`
}

type LikeError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type UnlikeSuccess struct {
	Status string `json:"status"`
}

type UnlikeError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type PostsSuccess struct {
	Status string `json:"status"`
	types.UsersPosts
	Message string `json:"message,omitempty"`
}

type PostsError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}
