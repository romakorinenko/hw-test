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
	DBPool, dbPoolErr := pgxpool.NewWithConfig(ctx, DBConfig)
	if dbPoolErr != nil {
		return nil, dbPoolErr
	}
	return DBPool, nil
}
