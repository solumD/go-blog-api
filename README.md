# go-blog-api
[Русский](#russian)

[English](#english)

[Эндпоинты (Endpoints)](#demo)
___
### RUS <a name="russian"></a> 
Небольшое API блога. Пользователи могут добавлять, редактировать и удалять посты. Также есть возможность увидеть посты конкретного пользователя и поставить лайк на конкретный пост.
В данный момент пишу функциональные тесты.
## Использованные пакеты и технологии
Роутер - [chi](https://github.com/go-chi/chi/);

Логгер - [slog](https://pkg.go.dev/golang.org/x/exp/slog);

Хранилище - [sqlite3](https://www.sqlite.org/) (при желании можно использовать любую базу данных, если реализовать используемые в хэндлерах интерфейсы);

Хэширование пароля - [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt);

JWT-аутентификация - [jwt-go](https://pkg.go.dev/github.com/golang-jwt/jwt/v5);

Документация - Swagger 2.0;

Написан DOCKERFILE;

Тестирование - стандартная библиотека и пакет [httpexpect](https://github.com/gavv/httpexpect).

___
### ENG <a name="english"></a> 
A small API of a blog. Users can add, update and delete posts. It is also possible to get the posts of a particular user and like a particular post.
Currently I'm writing functional tests.
## Packages and technologies used
Router - [chi](https://github.com/go-chi/chi/);

Logger - [slog](https://pkg.go.dev/golang.org/x/exp/slog);

Storage - [sqlite3](https://www.sqlite.org/) (it is possible to use any database you want, if you implement the interfaces used in the handlers);

Password hashing - [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt);

JWT authentication - [jwt-go](https://pkg.go.dev/github.com/golang-jwt/jwt/v5);

Documentation - Swagger 2.0;

DOCKERFILE is written;

Testing - standart library and [httpexpect](https://github.com/gavv/httpexpect).
___

## Эндпоинты (Endpoints) <a name="demo"></a> 

#### /auth/register - регистрация пользователя (user registration)

#### /auth/login - авторизация пользователя (user authorisation)

#### /user/{login} - получить все посты конкретного пользователя (get all posts of a particular user)

#### /post/create - создать пост (create a post)

#### /post/update - обновить название или текст поста (update post's title or text)

#### /post/delete - удалить пост (delete a post)

#### /post/like - поставить лайк на пост (like post)

#### /post/unlike - убрать лайк с поста (unlike post)
