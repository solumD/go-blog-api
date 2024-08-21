package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/solumD/go-blog-api/internal/http-server/handlers/post/remove"
	"github.com/solumD/go-blog-api/internal/http-server/handlers/post/save"
	"github.com/solumD/go-blog-api/internal/http-server/handlers/user/login"
	"github.com/solumD/go-blog-api/internal/http-server/handlers/user/register"
)

const (
	host = "localhost:8081"
)

func TestRegisterLoginCreateDelete(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	l := gofakeit.Username()
	p := gofakeit.Password(true, true, true, false, false, 10)

	e.POST("/auth/register").
		WithJSON(register.Request{
			Login:    l,
			Password: p,
		}).
		Expect().Status(http.StatusOK)

	r := e.POST("/auth/login").
		WithJSON(login.Request{
			Login:    l,
			Password: p,
		}).
		Expect().Status(http.StatusOK).
		JSON().Object().
		ContainsKey("token")

	token := r.Value("token").String().Raw()

	rr := e.POST("/post/create").
		WithJSON(save.Request{
			Title: gofakeit.Paragraph(1, 1, 3, ""),
			Text:  gofakeit.Paragraph(1, 5, 40, ""),
		}).
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK).
		JSON().Object().
		ContainsKey("id")

	postID := int(rr.Value("id").Number().Raw())

	e.DELETE("/post/delete").
		WithJSON(remove.Request{
			ID: postID,
		}).
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(200)

}

func TestRegisterUserAlreadyExists(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	l := gofakeit.Username()
	p := gofakeit.Password(true, true, true, false, false, 10)

	e.POST("/auth/register").
		WithJSON(register.Request{
			Login:    l,
			Password: p,
		}).
		Expect().Status(http.StatusOK)

	r := e.POST("/auth/register").
		WithJSON(register.Request{
			Login:    l,
			Password: p,
		}).
		Expect().Status(http.StatusBadRequest).
		JSON().Object().
		ContainsKey("error")

	expectedError := "user already exists"
	r.Value("error").String().IsEqual(expectedError)
}

// TODO: написать тест-кейсы
func TestRegisterInvalidLoginOrPassword(t *testing.T) {
	testCases := []struct {
		name     string
		login    string
		password string
		error    string
	}{
		{
			name:     "Login with space",
			login:    gofakeit.LetterN(3) + " " + gofakeit.Username(),
			password: gofakeit.Password(true, true, true, false, false, 10),
			error:    "login cannot contain spaces",
		},
		{
			name:     "Short login",
			login:    gofakeit.LetterN(7),
			password: gofakeit.Password(true, true, true, false, false, 10),
			error:    "login cannot be shorter than 8 characters",
		},
		{
			name:     "Password with space",
			login:    gofakeit.Username(),
			password: gofakeit.LetterN(3) + " " + gofakeit.Password(true, true, true, false, false, 10),
			error:    "password cannot contain spaces",
		},
		{
			name:     "Short password",
			login:    gofakeit.Username(),
			password: gofakeit.Password(true, true, true, false, false, 7),
			error:    "password cannot be shorter than 8 characters",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			r := e.POST("/auth/register").
				WithJSON(register.Request{
					Login:    tc.login,
					Password: tc.password,
				}).
				Expect().Status(http.StatusBadRequest).
				JSON().Object().
				ContainsKey("error")

			r.Value("error").String().IsEqual(tc.error)
		})
	}
}
