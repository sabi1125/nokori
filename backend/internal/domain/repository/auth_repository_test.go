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

func TestAuthRepository_InsertVerificationCode(t *testing.T) {
	testCode := entities.VerificationCodes{
		VerificationCodeId: "5fb7436d-b021-43cc-9dbf-231b2169164a",
		UserId:             "f79a7721-f7a4-441c-9136-a4db768c9bc3",
		Code:               "042917",
		ResendEmailId:      "9d950124-25f8-4b41-b88b-5aea038684f4",
		ExpiresAt:          time.Date(2026, 10, 10, 12, 1, 0, 0, time.UTC),
	}

	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		wantErrCode string
	}{
		{
			name: "inserts successfully",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `verification_codes`")).
					WithArgs(testCode.VerificationCodeId, testCode.UserId, testCode.Code,
						testCode.ResendEmailId, testCode.ExpiresAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "returns internal error when insert fails",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `verification_codes`")).
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

			repo := NewAuthRepository(db)
			err := repo.InsertVerificationCode(context.Background(), testCode)

			assertAppErrCode(t, tt.wantErrCode, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
