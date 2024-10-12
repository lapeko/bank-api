package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/random"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/repository"
	mockdb "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/repository/mocks"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TODO cover errors
func TestCreateUser(t *testing.T) {
	name := random.String(6)
	req := &createUserRequest{
		FullName: name,
		Email:    fmt.Sprintf("%s@mail.com", name),
		Password: name,
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), 4)
	require.NoError(t, err)
	user := repository.User{
		ID:                random.Int64(1, 1000),
		Email:             req.Email,
		FullName:          req.FullName,
		HashedPassword:    string(hashedPass),
		CreatedAt:         time.Now(),
		PasswordChangesAt: time.Now(),
	}
	cases := []struct {
		name          string
		request       *createUserRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{{name: "ok", request: req, buildStubs: func(store *mockdb.MockStore) {
		store.
			EXPECT().
			CreateUser(gomock.Any(), gomock.Any()).
			Times(1).
			Return(user, nil)
	}, checkResponse: func(recorder *httptest.ResponseRecorder) {
		fmt.Println(recorder)
		require.Equal(t, http.StatusCreated, recorder.Code)
		body, err := io.ReadAll(recorder.Body)
		require.NoError(t, err)
		var response struct {
			Ok   bool             `json:"ok"`
			Err  interface{}      `json:"err"`
			Body *repository.User `json:"body"`
		}
		err = json.Unmarshal(body, &response)
		require.NoError(t, err)
		require.True(t, response.Ok)
		require.Nil(t, response.Err)
		require.Equal(t, user.ID, response.Body.ID)
		require.Equal(t, user.Email, response.Body.Email)
		require.Equal(t, user.FullName, response.Body.FullName)
		require.Empty(t, response.Body.HashedPassword)
		require.WithinDuration(t, user.CreatedAt, response.Body.CreatedAt, time.Second)
		require.WithinDuration(t, user.PasswordChangesAt, response.Body.PasswordChangesAt, time.Second)
	}}}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			recorder := httptest.NewRecorder()

			path := fmt.Sprintf("/users/")
			body, err := json.Marshal(testCase.request)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(body))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")

			api := GetApi()
			api.store = store
			api.SetUpRoutes()

			api.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}
