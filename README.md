﻿# go-blog-api
[English](#english)
___
## RUS
Небольшое API блога. Пользователи могут добавлять и удалять посты. Также есть возможность увидеть посты конкретного пользователя.
## Использованные пакеты и технологии
Роутер - [chi](https://github.com/go-chi/chi/);

Логгер - [slog](https://pkg.go.dev/golang.org/x/exp/slog);

Хранилище - [sqlite3](https://www.sqlite.org/) (при желании можно использовать любую базу данных, если реализовать используемые в хэндлерах интерфейсы);

JWT-аутентификация;

Написан DOCKERFILE.

___
## ENG <a name="english"></a> 
A small API of a blog. Users can add and delete posts. It is also possible to get the posts of a particular user.
## Packages and technologies used
Router - [chi](https://github.com/go-chi/chi/);

Logger - [slog](https://pkg.go.dev/golang.org/x/exp/slog);

Storage - [sqlite3](https://www.sqlite.org/) (it is possible to use any database you want, if you implement the interfaces used in the handlers);

JWT authentication;

DOCKERFILE is written.
___

## Демонстрация эндпоинтов (Demonstration of endpoints)

#### /register - регистрация пользователя (user registration)

#### /login - авторизация пользователя (user authorisation)

#### /users/{login} - получить все посты конкретного пользователя (get all posts of a particular user)

#### /post/create - создать пост (create a post)

#### /post/delete - удалить пост (delete a post)
