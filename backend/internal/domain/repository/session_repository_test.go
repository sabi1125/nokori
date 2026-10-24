package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"backend/internal/domain/apperror"
	"backend/internal/domain/entities"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSessionRepository_CreateSession(t *testing.T) {
	testSession := entities.Sessions{
		SessionId: "0f8b2c6e-3d1a-4b7e-9c2f-5a6d7e8f9a0b",
		UserId:    "f79a7721-f7a4-441c-9136-a4db768c9bc3",
		TokenHash: "9f2c4e1b7a3d5f6e8c0b2a4d6f8e0c2b4a6d8f0e2c4b6a8d0f2e4c6b8a0d2f4e",
		ExpiresAt: time.Date(2026, 11, 21, 12, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		wantErrCode string
	}{
		{
			// previous_token_hash and previous_valid_until are nil, so they go in as NULL
			name: "inserts successfully",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `sessions`")).
					WithArgs(testSession.SessionId, testSession.UserId, testSession.TokenHash,
						nil, nil, testSession.ExpiresAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "returns internal error when insert fails",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `sessions`")).
					WillReturnError(errors.New("db connection lost"))
				mock.ExpectRollback()
			},
			wantErrCode: apperror.InternalError.Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			tt.setupMock(mock)

			repo := NewSessionRepository(db)
			err := repo.CreateSession(context.Background(), testSession)

			assertAppErrCode(t, tt.wantErrCode, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
