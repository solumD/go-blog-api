﻿# go-blog-api
Небольшое API блога. Пользователи могут добавлять и удалять посты. Также есть возможность посмотреть посты конкретного пользователя.
## Использованные пакеты и технологии
Роутер - [chi](https://github.com/go-chi/chi/);

Логгер - [slog](https://pkg.go.dev/golang.org/x/exp/slog);

Хранилище - [sqlite3](https://www.sqlite.org/) (при желании можно использовать любую базу данных, если реализовать используемые в хэндлерах интерфейсы);

JWT-аутентификация;

Написан DOCKERFILE для запуска в контейнере.

## Демонстрация эндпоинтов

### /login

### /register

### /users/{login}

### /post/create

### /post/delete
