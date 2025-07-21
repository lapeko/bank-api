package db

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	want := CreateUserParams{
		FullName:       utils.GenRandFullName(),
		Email:          utils.GetRandEmail(),
		HashedPassword: utils.GenRandHashedPassword(),
	}

	got, err := testStore.CreateUser(ctx, want)

	require.NoError(t, err)
	require.NotZero(t, got.ID)
	require.Equal(t, got.FullName, want.FullName)
	require.Equal(t, got.Email, want.Email)
	require.Equal(t, got.HashedPassword, want.HashedPassword)
	require.NotZero(t, got.CreatedAt)

	return got
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	got, err := testStore.GetUserByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.Equal(t, got, user)

	err = testStore.DeleteUser(ctx, user.ID)
	require.NoError(t, err)

	delUser, err := testStore.GetUserById(ctx, user.ID)
	require.Empty(t, delUser)
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestListUsers(t *testing.T) {
	defer cleanTestStore(t)

	params := ListUsersParams{Offset: 0, Limit: 2}
	got, err := testStore.ListUsers(ctx, params)
	require.NoError(t, err)
	require.Empty(t, got)

	var want []User
	for i := 0; i < 2; i++ {
		want = append(want, createRandomUser(t))
	}
	got, err = testStore.ListUsers(ctx, params)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestListUsers_QueryError(t *testing.T) {
	params := ListUsersParams{Offset: 0, Limit: 2}
	wantError := errors.New(utils.GenRandString(10))

	dbMock := new(dbConnMock)
	dbMock.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, wantError)

	gotUsers, gotError := New(dbMock).ListUsers(ctx, params)
	require.Nil(t, gotUsers)
	require.ErrorIs(t, gotError, wantError)
	dbMock.AssertCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestListUsers_ScanError(t *testing.T) {
	params := ListUsersParams{Offset: 0, Limit: 2}
	wantError := errors.New(utils.GenRandString(10))

	dbMock := new(dbConnMock)
	rowsMock := new(rowsMock)
	rowsMock.On("Next").Return(true)
	rowsMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(wantError)
	rowsMock.On("Close").Return()
	dbMock.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(rowsMock, nil)

	gotUsers, gotError := New(dbMock).ListUsers(ctx, params)
	require.Nil(t, gotUsers)
	require.ErrorIs(t, gotError, wantError)
	dbMock.AssertCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Next")
	rowsMock.AssertCalled(t, "Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Close")
}

func TestListUsers_RowsError(t *testing.T) {
	params := ListUsersParams{Offset: 0, Limit: 2}
	wantError := errors.New(utils.GenRandString(10))

	dbMock := new(dbConnMock)
	rowsMock := new(rowsMock)
	rowsMock.On("Next").Once().Return(true)
	rowsMock.On("Next").Once().Return(false)
	rowsMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	rowsMock.On("Err").Return(wantError)
	rowsMock.On("Close").Return()
	dbMock.On("Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(rowsMock, nil)

	gotUsers, gotError := New(dbMock).ListUsers(ctx, params)
	require.Nil(t, gotUsers)
	require.ErrorIs(t, gotError, wantError)
	dbMock.AssertCalled(t, "Query", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Next")
	rowsMock.AssertCalled(t, "Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rowsMock.AssertCalled(t, "Err")
	rowsMock.AssertCalled(t, "Close")
}

func TestUpdateUserEmail(t *testing.T) {
	defer cleanTestStore(t)

	user := createRandomUser(t)

	newEmail := utils.GetRandEmail()
	for newEmail == user.Email {
		newEmail = utils.GetRandEmail()
	}

	require.NotEqual(t, user.Email, newEmail)
	got, err := testStore.UpdateUserEmail(ctx, UpdateUserEmailParams{ID: user.ID, Email: newEmail})

	require.NoError(t, err)
	require.Equal(t, got.Email, newEmail)
}

func TestUpdateUserName(t *testing.T) {
	defer cleanTestStore(t)

	user := createRandomUser(t)

	newName := utils.GenRandFullName()
	for newName == user.FullName {
		newName = utils.GenRandFullName()
	}

	require.NotEqual(t, user.FullName, newName)
	got, err := testStore.UpdateUserFullName(ctx, UpdateUserFullNameParams{ID: user.ID, FullName: newName})

	require.NoError(t, err)
	require.Equal(t, got.FullName, newName)
}

func TestUpdateUserPassword(t *testing.T) {
	defer cleanTestStore(t)

	user := createRandomUser(t)

	newPassword := utils.GenRandHashedPassword()
	for newPassword == user.HashedPassword {
		newPassword = utils.GenRandHashedPassword()
	}

	require.NotEqual(t, user.HashedPassword, newPassword)
	got, err := testStore.UpdateUserPassword(ctx, UpdateUserPasswordParams{ID: user.ID, HashedPassword: newPassword})

	require.NoError(t, err)
	require.Equal(t, got.HashedPassword, newPassword)
}
