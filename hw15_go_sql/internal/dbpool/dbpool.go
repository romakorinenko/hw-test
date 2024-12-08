package dbpool

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romakorinenko/hw-test/hw15_go_sql/internal/config"
)

func NewDBPool(ctx context.Context, dBCfg *config.DB) (*pgxpool.Pool, error) {
	DBConfig, err := pgxpool.ParseConfig(dBCfg.ConnectionString)
	if err != nil {
		return nil, err
	}
	DBPool, err := pgxpool.NewWithConfig(ctx, DBConfig)
	if err != nil {
		return nil, err
	}
	return DBPool, nil
}
