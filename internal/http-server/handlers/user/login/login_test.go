package login_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/solumD/go-blog-api/internal/http-server/handlers/user/login"
	"github.com/solumD/go-blog-api/internal/http-server/handlers/user/login/mocks"
	"github.com/solumD/go-blog-api/internal/lib/logger/loggerdiscard"
	"github.com/solumD/go-blog-api/internal/lib/password"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLoginHandler(t *testing.T) {
	testcases := []struct {
		name                 string
		login                string
		password             string
		realPassword         string
		respError            string
		mockGetPasswordError error
		mockIsUserExistError error
		statusCode           int
	}{
		{
			name:         "Success",
			login:        "VeryCoolLogin",
			password:     "VeryCoolPassword",
			realPassword: "VeryCoolPassword",
			statusCode:   http.StatusOK,
		},
		{
			name:       "Short login",
			login:      "login",
			password:   "VeryCoolPassword",
			respError:  "login cannot be shorter than 8 characters",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Short password",
			login:      "VeryCoolLogin",
			password:   "pass",
			respError:  "password cannot be shorter than 8 characters",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Login with space",
			login:      "Very Cool Login",
			password:   "VeryCoolPassword",
			respError:  "login cannot contain spaces",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Password with space",
			login:      "VeryCoolLogin",
			password:   "Very Cool Password",
			respError:  "password cannot contain spaces",
			statusCode: http.StatusBadRequest,
		},
		{
			name:                 "GetPassword Error",
			login:                "VeryCoolLogin",
			password:             "VeryCoolPassword",
			respError:            "failed to get user's real password",
			mockGetPasswordError: errors.New("unexpected error"),
			statusCode:           http.StatusInternalServerError,
		},
		{
			name:                 "IsUserExist Error",
			login:                "VeryCoolLogin",
			password:             "VeryCoolPassword",
			respError:            "failed to check if user exists",
			mockIsUserExistError: errors.New("unexpected error"),
			statusCode:           http.StatusInternalServerError,
		},
		{
			name:       "User does not exist",
			login:      "VeryCoolLogin",
			password:   "VeryCoolPassword",
			respError:  "user does not exist",
			statusCode: http.StatusBadRequest,
		},
		{
			name:         "Passwords doesn't match",
			login:        "VeryCoolLogin",
			password:     "VeryCoolPassword",
			realPassword: "AnotherPassword",
			respError:    "invalid password",
			statusCode:   http.StatusBadRequest,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userAuthorizerMock := mocks.NewUserAuthorizer(t)

			hashedrealPassword, err := password.EncryptPassword(tc.realPassword)
			require.NoError(t, err)

			if tc.respError == "" || tc.respError == "invalid password" || tc.mockGetPasswordError != nil {
				userAuthorizerMock.On("GetPassword", mock.Anything, mock.AnythingOfType("string")).
					Return(hashedrealPassword, tc.mockGetPasswordError).
					Once()
				userAuthorizerMock.On("IsUserExist", mock.Anything, mock.AnythingOfType("string")).
					Return(true, nil).
					Once()
			}

			if tc.respError == "user does not exist" || tc.mockIsUserExistError != nil {
				userAuthorizerMock.On("IsUserExist", mock.Anything, mock.AnythingOfType("string")).
					Return(false, tc.mockIsUserExistError).
					Once()
			}

			handler := login.New(context.Background(), "secret", loggerdiscard.NewDiscardLogger(), userAuthorizerMock)

			input := fmt.Sprintf(`{"login": "%s", "password": "%s"}`, tc.login, tc.password)

			req, err := http.NewRequest(http.MethodPost, "/user/login", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			require.Equal(t, recorder.Code, tc.statusCode)

			body := recorder.Body.String()

			var resp login.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
