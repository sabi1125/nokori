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
						testCode.ResendEmailId, testCode.Attempts, testCode.ExpiresAt).
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

func TestAuthRepository_GetVerifiationCodeByEmail(t *testing.T) {
	const userId = "f79a7721-f7a4-441c-9136-a4db768c9bc3"
	expiresAt := time.Date(2026, 10, 10, 12, 1, 0, 0, time.UTC)
	columns := []string{"verification_code_id", "user_id", "code", "resend_email_id", "attempts", "expires_at"}

	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		wantCode    entities.VerificationCodes
		wantErrCode string
	}{
		{
			name: "returns the code for the user",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `verification_codes` WHERE user_id = ?")).
					WithArgs(userId, 1).
					WillReturnRows(sqlmock.NewRows(columns).
						AddRow("code-1", userId, "042917", "resend-1", 2, expiresAt))
			},
			wantCode: entities.VerificationCodes{
				VerificationCodeId: "code-1",
				UserId:             userId,
				Code:               "042917",
				ResendEmailId:      "resend-1",
				Attempts:           2,
				ExpiresAt:          expiresAt,
			},
		},
		{
			// signup always inserts a code with the user, so a missing one is unexpected
			name: "returns internal error when the user has no code",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `verification_codes` WHERE user_id = ?")).
					WithArgs(userId, 1).
					WillReturnRows(sqlmock.NewRows(columns))
			},
			wantErrCode: apperror.InternalError.Code,
		},
		{
			name: "returns internal error when the query fails",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `verification_codes` WHERE user_id = ?")).
					WithArgs(userId, 1).
					WillReturnError(errors.New("db connection lost"))
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
			code, err := repo.GetVerifiationCodeByEmail(context.Background(), userId)

			assertAppErrCode(t, tt.wantErrCode, err)
			assert.Equal(t, tt.wantCode, code)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthRepository_UpdateVerificationCodeAttempts(t *testing.T) {
	const verificationCodeId = "5fb7436d-b021-43cc-9dbf-231b2169164a"

	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		wantErrCode string
	}{
		{
			name: "updates successfully",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `verification_codes` SET `attempts`=? WHERE verification_code_id = ?")).
					WithArgs(3, verificationCodeId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "returns internal error when update fails",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE `verification_codes` SET `attempts`=? WHERE verification_code_id = ?")).
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
			err := repo.UpdateVerificationCodeAttempts(context.Background(), verificationCodeId, 3)

			assertAppErrCode(t, tt.wantErrCode, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
