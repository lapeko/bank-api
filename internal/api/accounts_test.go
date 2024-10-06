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
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountApi(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	store.
		EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	recorder := httptest.NewRecorder()

	path := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)

	api := New()
	api.store = store
	api.SetUpRoutes()

	api.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	matchAccountBody(t, account, recorder.Body)
}

func randomAccount() repository.Account {
	return repository.Account{
		ID:       random.Int64(1, 1000),
		Owner:    random.String(5),
		Balance:  random.Int64(0, 1000),
		Currency: random.Currency(),
	}
}

func matchAccountBody(t *testing.T, account repository.Account, body *bytes.Buffer) {
	var response struct {
		Ok   bool               `json:"ok"`
		Err  interface{}        `json:"err"`
		Body repository.Account `json:"body"`
	}
	err := json.NewDecoder(body).Decode(&response)
	require.NoError(t, err)
	require.Equal(t, account, response.Body)
}
