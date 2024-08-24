package save_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/solumD/go-blog-api/internal/http-server/handlers/post/save"
	"github.com/solumD/go-blog-api/internal/http-server/handlers/post/save/mocks"
	"github.com/solumD/go-blog-api/internal/lib/logger/loggerdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSavePostHandler(t *testing.T) {
	testCases := []struct {
		name       string
		title      string
		text       string
		respError  string
		mockError  error
		statusCode int
	}{
		{
			name:       "Success",
			title:      "Very Cool Title",
			text:       "Very Cool Text",
			statusCode: http.StatusOK,
		},
		{
			name:       "Empty title",
			title:      "",
			text:       "Very Cool Text",
			respError:  "post's title and text can't be empty",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Empty text",
			title:      "Very Cool Title",
			text:       "",
			respError:  "post's title and text can't be empty",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "SavePost Error",
			title:      "Very Cool Title",
			text:       "Very Cool Text",
			respError:  "failed to save post",
			mockError:  errors.New("unexpected error"),
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			postSaverMock := mocks.NewPostSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				postSaverMock.On("SavePost", mock.Anything, mock.AnythingOfType("string"), tc.title, tc.text, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := save.New(context.Background(), loggerdiscard.NewDiscardLogger(), postSaverMock)

			input := fmt.Sprintf(`{"title": "%s", "text": "%s"}`, tc.title, tc.text)

			req, err := http.NewRequest(http.MethodPost, "/post/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			req.Header.Add("login", "test_user")

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			require.Equal(t, recorder.Code, tc.statusCode)

			body := recorder.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
