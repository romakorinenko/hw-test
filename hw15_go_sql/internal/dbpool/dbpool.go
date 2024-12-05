package dbpool

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romakorinenko/hw-test/hw15_go_sql/internal/config"
)

func NewDbPool(ctx context.Context, dbCfg *config.Db) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(dbCfg.ConnectionString)
	if err != nil {
		return nil, err
	}
	dbPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}
	return dbPool, nil
}
