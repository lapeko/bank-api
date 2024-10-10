package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/random"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/repository"
	mockdb "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/storage/repository/mocks"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountApi(t *testing.T) {
	account := randomAccount()

	cases := []struct {
		name          string
		accountId     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "ok",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				matchAccountResponse(t, true, nil, account, recorder.Body)
			},
		},
		{
			name:      "not found",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(repository.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				matchAccountResponse(t, false, "user not found", nil, recorder.Body)
			},
		},
		{
			name:      "db error",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(repository.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				matchAccountResponse(t, false, "sql: connection is already closed", nil, recorder.Body)
			},
		},
		{
			name:      "path param parse error",
			accountId: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				matchAccountResponse(
					t,
					false,
					"Key: 'getAccountByIdRequest.Id' Error:Field validation for 'Id' failed on the 'required' tag",
					nil,
					recorder.Body,
				)
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			recorder := httptest.NewRecorder()

			path := fmt.Sprintf("/accounts/%d", testCase.accountId)
			request, err := http.NewRequest(http.MethodGet, path, nil)
			require.NoError(t, err)

			api := GetApi()
			api.store = store
			api.SetUpRoutes()

			api.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func randomAccount() repository.Account {
	return repository.Account{
		ID:       random.Int64(1, 1000),
		UserID:   random.Int64(1, 1000),
		Balance:  random.Int64(0, 1000),
		Currency: random.Currency(),
	}
}

func matchAccountResponse(t *testing.T, ok bool, err interface{}, body interface{}, httpResponse *bytes.Buffer) {
	var response struct {
		Ok   bool                `json:"ok"`
		Err  interface{}         `json:"err"`
		Body *repository.Account `json:"body"`
	}
	e := json.NewDecoder(httpResponse).Decode(&response)
	require.NoError(t, e)
	require.Equal(t, ok, response.Ok)
	require.Equal(t, err, response.Err)
	if response.Body == nil {
		require.Nil(t, body)
		return
	}
	require.Equal(t, body, *response.Body)
}
