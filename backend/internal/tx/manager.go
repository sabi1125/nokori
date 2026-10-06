//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
package tx

import (
	"context"

	"gorm.io/gorm"
)

type Manager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type ctxKey string

const gormTxKey ctxKey = "gorm_tx"

func WithTx(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, gormTxKey, db)
}

func ExtractTx(ctx context.Context) *gorm.DB {
	db, _ := ctx.Value(gormTxKey).(*gorm.DB)
	return db
}
