package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type dbConnMock struct {
	mock.Mock
}

func (m *dbConnMock) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *dbConnMock) Commit(ctx context.Context) error {
	callArgs := m.Called(ctx)
	return callArgs.Error(0)
}

func (m *dbConnMock) Rollback(ctx context.Context) error {
	callArgs := m.Called(ctx)
	return callArgs.Error(0)
}

func (m *dbConnMock) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	callArgs := m.Called(ctx, tableName, columnNames, rowSrc)
	return callArgs.Get(0).(int64), callArgs.Error(1)
}

func (m *dbConnMock) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	callArgs := m.Called(ctx, b)
	return callArgs.Get(0).(pgx.BatchResults)
}

func (m *dbConnMock) LargeObjects() pgx.LargeObjects {
	callArgs := m.Called()
	return callArgs.Get(0).(pgx.LargeObjects)
}

func (m *dbConnMock) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	callArgs := m.Called(ctx, name, sql)
	return callArgs.Get(0).(*pgconn.StatementDescription), callArgs.Error(1)
}

func (m *dbConnMock) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	callArgs := m.Called(append([]interface{}{ctx, query}, args...)...)
	cTag, _ := callArgs.Get(0).(pgconn.CommandTag)
	err, _ := callArgs.Get(1).(error)
	return cTag, err
}

func (m *dbConnMock) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	callArgs := m.Called(append([]interface{}{ctx, query}, args...)...)
	rows, _ := callArgs.Get(0).(pgx.Rows)
	err, _ := callArgs.Get(1).(error)
	return rows, err
}

func (m *dbConnMock) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	callArgs := m.Called(append([]interface{}{ctx, query}, args...)...)
	row, _ := callArgs.Get(0).(pgx.Row)
	return row
}

func (m *dbConnMock) Conn() *pgx.Conn {
	return m.Called().Get(0).(*pgx.Conn)
}

func (m *dbConnMock) GetAccountsByIdForUpdate(ctx context.Context, arg GetAccountsByIdForUpdateParams) ([]Account, error) {
	callArgs := m.Called(ctx, arg)
	return callArgs.Get(0).([]Account), callArgs.Error(1)
}

type rowsMock struct {
	mock.Mock
}

func (m *rowsMock) Close() {
	m.Called()
}

func (m *rowsMock) Err() error {
	return m.Called().Error(0)
}

func (m *rowsMock) CommandTag() pgconn.CommandTag {
	return m.Called().Get(0).(pgconn.CommandTag)
}

func (m *rowsMock) FieldDescriptions() []pgconn.FieldDescription {
	return m.Called().Get(0).([]pgconn.FieldDescription)
}

func (m *rowsMock) Next() bool {
	return m.Called().Bool(0)
}

func (m *rowsMock) Scan(dest ...any) error {
	return m.Called(dest...).Error(0)
}

func (m *rowsMock) Values() ([]any, error) {
	args := m.Called()
	return args.Get(0).([]any), args.Error(1)
}

func (m *rowsMock) RawValues() [][]byte {
	return m.Called().Get(0).([][]byte)
}

func (m *rowsMock) Conn() *pgx.Conn {
	return m.Called().Get(0).(*pgx.Conn)
}
