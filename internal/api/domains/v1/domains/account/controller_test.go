package account

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	apiUtils "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/api/utils"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/mockdb"
	db "github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils"
	"github.com/stretchr/testify/require"
)

var path = struct {
	accounts string
	v1       string
}{accounts: "/accounts", v1: "/v1"}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	apiUtils.RegisterValidators()
	os.Exit(m.Run())
}

func TestCreateAccountHandler(t *testing.T) {
	userID := int64(1)
	currency := utils.CurrencyUSD
	account := db.Account{
		ID:       1,
		UserID:   userID,
		Currency: currency,
		Balance:  0,
	}

	tests := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(*httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"user_id":  userID,
				"currency": currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					UserID:   userID,
					Currency: currency,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, rr.Code)
			},
		},
		{
			name: "InvalidBody",
			body: gin.H{
				"user_id": 0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rr.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"user_id":  userID,
				"currency": currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					UserID:   userID,
					Currency: currency,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), arg).
					Times(1).
					Return(db.Account{}, errors.New("db error"))
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rr.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tt.buildStubs(store)

			router := gin.Default()
			Register(path.accounts, router.Group(path.v1), store)

			rec := httptest.NewRecorder()
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s%s", path.v1, path.accounts, "/"), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rec, req)
			tt.checkResponse(rec)
		})
	}
}

func TestListAccountsHandler(t *testing.T) {
	tests := []struct {
		name          string
		queryParams   string
		buildStubs    func(*mockdb.MockStore)
		checkResponse func(*httptest.ResponseRecorder)
	}{
		{
			name:        "ListAccountsHandler should send error if ShouldBind failed because of no query params",
			queryParams: "",
			buildStubs:  func(ms *mockdb.MockStore) {},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, rr.Code, http.StatusBadRequest)
				var got gin.H
				err := json.Unmarshal(rr.Body.Bytes(), &got)
				require.NoError(t, err)
				require.NotEmpty(t, got["error"])
			},
		},
		{
			name:        "ListAccountsHandler should send error if store.ListAccounts returned pgx.ErrNoRows",
			queryParams: "?page=1&size=5",
			buildStubs: func(ms *mockdb.MockStore) {
				ms.
					EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(db.ListAccountsParams{Limit: 5, Offset: 0})).
					Times(1).
					Return([]db.ListAccountsRow{}, pgx.ErrNoRows)
				ms.
					EXPECT().
					GetTotalAccountsCount(gomock.Any()).
					Times(1).
					Return(int64(0), nil)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rr.Code)
				var got gin.H
				err := json.Unmarshal(rr.Body.Bytes(), &got)
				require.NoError(t, err)
				require.Equal(t, got["error"], pgx.ErrNoRows.Error())
			},
		},
		{
			name:        "ListAccountsHandler should send error if store.ListAccounts failed not with pgx.ErrNoRows",
			queryParams: "?page=1&size=5",
			buildStubs: func(ms *mockdb.MockStore) {
				ms.
					EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(db.ListAccountsParams{Limit: 5, Offset: 0})).
					Times(1).
					Return([]db.ListAccountsRow{{ID: int64(1)}}, errors.New("Some random error"))
				ms.
					EXPECT().
					GetTotalAccountsCount(gomock.Any()).
					Times(1).
					Return(int64(0), nil)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rr.Code)
				var got gin.H
				err := json.Unmarshal(rr.Body.Bytes(), &got)
				require.NoError(t, err)
				require.Equal(t, got["error"], "Some random error")
			},
		},
		{
			name:        "ListAccountsHandler should success",
			queryParams: "?page=2&size=10",
			buildStubs: func(ms *mockdb.MockStore) {
				ms.
					EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(db.ListAccountsParams{Limit: 10, Offset: 10})).
					Times(1).
					Return([]db.ListAccountsRow{{ID: int64(123)}}, nil)
				ms.
					EXPECT().
					GetTotalAccountsCount(gomock.Any()).
					Times(1).
					Return(int64(20), nil)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rr.Code)
				var got gin.H
				err := json.Unmarshal(rr.Body.Bytes(), &got)
				require.NoError(t, err)
				require.Equal(t, got["error"], "")
				require.Equal(t, got["ok"], true)

				dataJson, err := json.Marshal(got["data"])
				require.NoError(t, err)

				var data listAccountsResponse
				json.Unmarshal(dataJson, &data)

				require.Equal(t, int64(20), data.TotalCount)
				require.Equal(t, 1, len(data.Accounts))
				require.Equal(t, int64(123), data.Accounts[0].ID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			router := gin.Default()
			store := mockdb.NewMockStore(ctrl)

			Register(path.accounts, router.Group(path.v1), store)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s/%s", path.v1, path.accounts, tt.queryParams), bytes.NewReader([]byte{}))

			tt.buildStubs(store)

			router.ServeHTTP(rec, req)
			tt.checkResponse(rec)
		})
	}
}

func TestGetAccountByIdHandler(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		buildStubs    func(*mockdb.MockStore)
		checkResponse func(*httptest.ResponseRecorder)
	}{
		{
			name: "Should bind URI fails",
			url:  "/v1/accounts/invalid-id",
			buildStubs: func(ms *mockdb.MockStore) {
				ms.EXPECT().GetAccountById(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rr.Code)
			},
		},
		{
			name: "Account not found",
			url:  "/v1/accounts/1",
			buildStubs: func(ms *mockdb.MockStore) {
				ms.EXPECT().
					GetAccountById(gomock.Any(), int64(1)).
					Times(1).
					Return(db.GetAccountByIdRow{}, pgx.ErrNoRows)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rr.Code)
			},
		},
		{
			name: "Internal error",
			url:  "/v1/accounts/1",
			buildStubs: func(ms *mockdb.MockStore) {
				ms.EXPECT().
					GetAccountById(gomock.Any(), int64(1)).
					Times(1).
					Return(db.GetAccountByIdRow{}, errors.New("db error"))
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rr.Code)
			},
		},
		{
			name: "Success",
			url:  "/v1/accounts/1",
			buildStubs: func(ms *mockdb.MockStore) {
				ms.EXPECT().
					GetAccountById(gomock.Any(), int64(1)).
					Times(1).
					Return(db.GetAccountByIdRow{ID: int64(23)}, nil)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rr.Code)
				var got gin.H
				err := json.Unmarshal(rr.Body.Bytes(), &got)
				require.NoError(t, err)
				require.Equal(t, true, got["ok"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tt.buildStubs(store)

			router := gin.Default()
			Register(path.accounts, router.Group(path.v1), store)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)

			router.ServeHTTP(rec, req)
			tt.checkResponse(rec)
		})
	}
}

func TestDeleteAccountHandler(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		buildStubs    func(*mockdb.MockStore)
		checkResponse func(*httptest.ResponseRecorder)
	}{
		{
			name: "Should bind URI fails",
			url:  "/v1/accounts/invalid-id",
			buildStubs: func(ms *mockdb.MockStore) {
				ms.EXPECT().DeleteAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rr.Code)
			},
		},
		{
			name: "Account not found",
			url:  "/v1/accounts/1",
			buildStubs: func(ms *mockdb.MockStore) {
				ms.EXPECT().
					DeleteAccount(gomock.Any(), int64(1)).
					Times(1).
					Return(db.Account{}, pgx.ErrNoRows)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rr.Code)
			},
		},
		{
			name: "Internal error",
			url:  "/v1/accounts/1",
			buildStubs: func(ms *mockdb.MockStore) {
				ms.EXPECT().
					DeleteAccount(gomock.Any(), int64(1)).
					Times(1).
					Return(db.Account{}, errors.New("db error"))
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rr.Code)
			},
		},
		{
			name: "Success",
			url:  "/v1/accounts/1",
			buildStubs: func(ms *mockdb.MockStore) {
				ms.EXPECT().
					DeleteAccount(gomock.Any(), int64(1)).
					Times(1).
					Return(db.Account{ID: int64(25)}, nil)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rr.Code)
				var got gin.H
				err := json.Unmarshal(rr.Body.Bytes(), &got)
				require.NoError(t, err)
				require.Equal(t, true, got["ok"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tt.buildStubs(store)

			router := gin.Default()
			Register(path.accounts, router.Group(path.v1), store)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, tt.url, nil)

			router.ServeHTTP(rec, req)
			tt.checkResponse(rec)
		})
	}
}
