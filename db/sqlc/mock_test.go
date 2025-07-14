package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type DBTXMock struct {
	mock.Mock
}

func (m *DBTXMock) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	callArgs := m.Called(append([]interface{}{ctx, query}, args...)...)
	cTag, _ := callArgs.Get(0).(pgconn.CommandTag)
	err, _ := callArgs.Get(1).(error)
	return cTag, err
}

func (m *DBTXMock) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	callArgs := m.Called(append([]interface{}{ctx, query}, args...)...)
	rows, _ := callArgs.Get(0).(pgx.Rows)
	err, _ := callArgs.Get(1).(error)
	return rows, err
}

func (m *DBTXMock) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	callArgs := m.Called(append([]interface{}{ctx, query}, args...)...)
	row, _ := callArgs.Get(0).(pgx.Row)
	return row
}

type RowsMock struct {
	mock.Mock
}

func (m *RowsMock) Close() {
	m.Called()
}

func (m *RowsMock) Err() error {
	return m.Called().Error(0)
}

func (m *RowsMock) CommandTag() pgconn.CommandTag {
	return m.Called().Get(0).(pgconn.CommandTag)
}

func (m *RowsMock) FieldDescriptions() []pgconn.FieldDescription {
	return m.Called().Get(0).([]pgconn.FieldDescription)
}

func (m *RowsMock) Next() bool {
	return m.Called().Bool(0)
}

func (m *RowsMock) Scan(dest ...any) error {
	return m.Called(dest...).Error(0)
}

func (m *RowsMock) Values() ([]any, error) {
	args := m.Called()
	return args.Get(0).([]any), args.Error(1)
}

func (m *RowsMock) RawValues() [][]byte {
	return m.Called().Get(0).([][]byte)
}

func (m *RowsMock) Conn() *pgx.Conn {
	return m.Called().Get(0).(*pgx.Conn)
}
