# go-blog-api
[Русский](#russian)

[English](#english)

[Эндпоинты (Endpoints)](#demo)
___
### RUS <a name="russian"></a> 
Небольшое API блога. Пользователи могут добавлять, редактировать и удалять посты. Также есть возможность увидеть посты конкретного пользователя и поставить лайк на конкретный пост.
## Использованные пакеты и технологии
Роутер - [chi](https://github.com/go-chi/chi/);

Логгер - [slog](https://pkg.go.dev/golang.org/x/exp/slog);

Хранилище - [sqlite3](https://www.sqlite.org/) (при желании можно использовать любую базу данных, если реализовать используемые в хэндлерах интерфейсы);

Хэширование пароля - [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt);

JWT-аутентификация - [jwt-go](https://pkg.go.dev/github.com/golang-jwt/jwt/v5);

Документация - Swagger 2.0;

Написан Dockerfile;

Моки для тестов - [mockery](https://github.com/vektra/mockery);

Тестирование - стандартная библиотека и пакет [httpexpect](https://github.com/gavv/httpexpect).

___
### ENG <a name="english"></a> 
A small API of a blog. Users can add, update and delete posts. It is also possible to get the posts of a particular user and like a particular post.
## Packages and technologies used
Router - [chi](https://github.com/go-chi/chi/);

Logger - [slog](https://pkg.go.dev/golang.org/x/exp/slog);

Storage - [sqlite3](https://www.sqlite.org/) (it is possible to use any database you want, if you implement the interfaces used in the handlers);

Password hashing - [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt);

JWT authentication - [jwt-go](https://pkg.go.dev/github.com/golang-jwt/jwt/v5);

Documentation - Swagger 2.0;

Dockerfile is written;

Mocks - [mockery](https://github.com/vektra/mockery);

Testing - standart library and [httpexpect](https://github.com/gavv/httpexpect).
___

## Эндпоинты (Endpoints) <a name="demo"></a> 

#### POST /auth/register - регистрация пользователя (user registration)
 
##### Example Input: 
```
{
    "login": "cool_user",
    "password": "qwerty123"
}
```

##### Example Response: 
```
{
    "status": "OK",
    "id": 7
} 
```

#### POST /auth/login - авторизация пользователя (user authorisation)

##### Example Input: 
```
{
    "login": "cool_user",
    "password": "qwerty123"
}
```

##### Example Response: 
```
{
    "status": "OK",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjI3ODA4OTIsInN1YiI6ImNvb2xfdXNlciJ9.n-c0YymZyImpDRCui4WLKmH25mzhdX2q_kTYyAIHpXE"
}
```

#### GET /user/{login} - получить все посты конкретного пользователя (get all posts of a particular user)

##### Example Response: 
```
{
    "status": "OK",
    "posts": [
        {
            "id": 3,
            "created_by": "prince123",
            "title": "aaaaaaaa",
            "text": "3333333333",
            "likes": 2,
            "liked_by": [
                "prince123",
                "solum123"
            ],
            "created_at": "2024-08-01T17:18:30Z",
            "updated_at": "2024-08-01T17:20:49Z"
        },
        {
            "id": 1,
            "created_by": "prince123",
            "title": "post1",
            "text": "aaaaaaaa",
            "likes": 2,
            "liked_by": [
                "prince123",
                "testuser"
            ],
            "created_at": "2024-08-01T17:18:21Z",
            "updated_at": "2024-08-01T17:20:55Z"
        }
    ]
}
```

#### POST /post/create - создать пост (create a post)

##### Example Input: 
```
{
    "title": "very cool title",
    "text": "very cool text"
}
```

##### Example Response: 
```
{
    "status": "OK",
    "id": 6
}
```

#### PATCH /post/update - обновить название или текст поста (update post's title or text)

##### Example Input: 
```
{   
    "id": 6,
    "title": "new title",
    "text": "new text"
}
```

##### Example Response: 
```
{
    "status": "OK"
}
```

#### DELETE /post/delete - удалить пост (delete a post)

##### Example Input: 
```
{   
    "id": 6
}
```

##### Example Response: 
```
{
    "status": "OK"
}
```

#### PUT /post/like - поставить лайк на пост (like a post)

##### Example Input: 
```
{   
    "id": 3
}
```

##### Example Response: 
```
{
    "status": "OK"
}
```

#### PUT /post/unlike - убрать лайк с поста (unlike a post)

##### Example Input: 
```
{   
    "id": 3
}
```

##### Example Response: 
```
{
    "status": "OK"
}
```
