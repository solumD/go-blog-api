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

	e.POST("/register").
		WithJSON(register.Request{
			Login:            l,
			Password:         p,
			RepeatedPassword: p,
		}).
		Expect().Status(http.StatusOK)

	r := e.POST("/login").
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
