package register_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/solumD/go-blog-api/internal/http-server/handlers/user/register"
	"github.com/solumD/go-blog-api/internal/http-server/handlers/user/register/mocks"
	"github.com/solumD/go-blog-api/internal/lib/logger/loggerdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegisterHandler(t *testing.T) {
	testcases := []struct {
		name                 string
		login                string
		password             string
		respError            string
		mockSaveError        error
		mockIsUserExistError error
		statusCode           int
	}{
		{
			name:       "Success",
			login:      "VeryCoolLogin",
			password:   "VeryCoolPassword",
			statusCode: http.StatusOK,
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
			name:          "SaveUser Error",
			login:         "VeryCoolLogin",
			password:      "VeryCoolPassword",
			respError:     "failed to save user",
			mockSaveError: errors.New("unexpected error"),
			statusCode:    http.StatusInternalServerError,
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
			name:       "User already exists",
			login:      "VeryCoolLogin",
			password:   "VeryCoolPassword",
			respError:  "user already exists",
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userRegistrarMock := mocks.NewUserRegistrar(t)

			if tc.respError == "" || tc.mockSaveError != nil {
				userRegistrarMock.On("SaveUser", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(int64(1), tc.mockSaveError).
					Once()
				userRegistrarMock.On("IsUserExist", mock.Anything, mock.AnythingOfType("string")).
					Return(false, nil).
					Once()
			}

			if tc.respError == "user already exists" || tc.mockIsUserExistError != nil {
				userRegistrarMock.On("IsUserExist", mock.Anything, mock.AnythingOfType("string")).
					Return(true, tc.mockIsUserExistError).
					Once()
			}

			handler := register.New(context.Background(), loggerdiscard.NewDiscardLogger(), userRegistrarMock)

			input := fmt.Sprintf(`{"login": "%s", "password": "%s"}`, tc.login, tc.password)

			req, err := http.NewRequest(http.MethodPost, "/user/login", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			require.Equal(t, recorder.Code, tc.statusCode)

			body := recorder.Body.String()

			var resp register.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
