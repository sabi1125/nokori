package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"backend/internal/domain/apperror"
	"backend/internal/domain/entities"
	"backend/internal/tx"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserRepository_CreateUser(t *testing.T) {
	testUser := entities.Users{
		UserID:       "01a116b7-9857-72d0-8afc-4cd265f71b25",
		FirstName:    "sabir",
		LastName:     "barahi",
		Email:        "example@gmail.com",
		PasswordHash: "$2a$10$hash",
		Verified:     false,
	}

	tests := []struct {
		name        string
		inTx        bool
		setupMock   func(mock sqlmock.Sqlmock)
		wantErrCode string
	}{
		{
			name: "creates successfully",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "returns internal error when insert fails",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WillReturnError(errors.New("db connection lost"))
				mock.ExpectRollback()
			},
			wantErrCode: apperror.InternalError.Code,
		},
		{
			// Signup inserts the user and the verification code in one
			// transaction, so the insert must run on the tx from ctx instead
			// of opening its own.
			name: "uses the transaction from context",
			inTx: true,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			tt.setupMock(mock)

			// In the tx case the insert must follow our Begin directly; if the
			// repository ignored the tx it would Begin a second time and
			// sqlmock would fail on the unexpected call.
			ctx := context.Background()
			var gormTx *gorm.DB
			if tt.inTx {
				gormTx = db.Begin()
				ctx = tx.WithTx(ctx, gormTx)
			}

			repo := NewUserRepository(db)
			err := repo.CreateUser(ctx, testUser)

			if tt.inTx {
				gormTx.Commit()
			}

			assertAppErrCode(t, tt.wantErrCode, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_FindUserByEmail(t *testing.T) {
	const email = "example@gmail.com"
	columns := []string{"user_id", "first_name", "last_name", "email", "password_hash", "verified"}

	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		wantUser    entities.Users
		wantErrCode string
	}{
		{
			name: "returns the user when the email exists",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email= ?")).
					WithArgs(email, 1).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow("user-1", "sabir", "barahi", email, "$2a$10$hash", false))
			},
			wantUser: entities.Users{
				UserID:       "user-1",
				FirstName:    "sabir",
				LastName:     "barahi",
				Email:        email,
				PasswordHash: "$2a$10$hash",
				Verified:     false,
			},
		},
		{
			// Not found is the normal signup case (email is free), so it must
			// come back as an empty user and no error.
			name: "returns an empty user and no error when the email does not exist",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email= ?")).
					WithArgs(email, 1).
					WillReturnRows(sqlmock.NewRows(columns))
			},
			wantUser: entities.Users{},
		},
		{
			name: "returns internal error when the query fails",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email= ?")).
					WithArgs(email, 1).
					WillReturnError(errors.New("db connection lost"))
			},
			wantUser:    entities.Users{},
			wantErrCode: apperror.InternalError.Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			tt.setupMock(mock)

			repo := NewUserRepository(db)
			user, err := repo.FindUserByEmail(context.Background(), email)

			assertAppErrCode(t, tt.wantErrCode, err)
			assert.Equal(t, tt.wantUser, user)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// assertAppErrCode checks err is nil when wantCode is empty, otherwise that
// it's an AppError with that code.
func assertAppErrCode(t *testing.T, wantCode string, err error) {
	t.Helper()
	if wantCode == "" {
		assert.NoError(t, err)
		return
	}
	appErr, ok := errors.AsType[*apperror.AppError](err)
	if assert.True(t, ok, "expected *apperror.AppError, got %T", err) {
		assert.Equal(t, wantCode, appErr.Code)
	}
}
