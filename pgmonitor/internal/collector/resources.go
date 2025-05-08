package collector

import (
	"context"
	"database/sql"
	"time"

	"pgmonitor/internal/models"
	"pgmonitor/internal/repository"
	"pgmonitor/pkg/logger"
	"pgmonitor/config"
)

type ResourceCollector struct {
	repo   *repository.Repository
	logger *logger.Logger
	config *config.Config
}

func NewResourceCollector(repo *repository.Repository, logger *logger.Logger, cfg *config.Config) *ResourceCollector {
	return &ResourceCollector{
		repo:   repo,
		logger: logger,
		config: cfg,
	}
}

func (c *ResourceCollector) Collect(ctx context.Context, sourceDB *sql.DB) error {
	c.logger.Info("Collecting resource usage statistics")
	
	// Получаем время старта сервера
	var pgStartTime time.Time
	err := sourceDB.QueryRowContext(ctx, "SELECT pg_postmaster_start_time()").Scan(&pgStartTime)
	if err != nil {
		return err
	}
	uptime := time.Since(pgStartTime).Seconds()
	
	// Получаем max_connections
	var maxConnections int
	err = sourceDB.QueryRowContext(ctx, 
		"SELECT setting FROM pg_settings WHERE name = 'max_connections'").Scan(&maxConnections)
	if err != nil {
		return err
	}
	
	query := `
		WITH db_stats AS (
			SELECT datname, COUNT(*) AS active_connections FROM pg_stat_activity GROUP BY datname
		),
		cache_stats AS (
			SELECT datname, sum(blks_hit)/nullif(sum(blks_hit+blks_read),0)*100 AS cache_hit_ratio 
			FROM pg_stat_database GROUP BY datname
		),
		txn_stats AS (
			SELECT datname, xact_commit+xact_rollback AS total_transactions,
				   tup_fetched, tup_inserted, tup_updated, tup_deleted
			FROM pg_stat_database
		)
		SELECT
			db.datname,
			db.active_connections,
			c.cache_hit_ratio,
			t.total_transactions,
			t.tup_fetched,
			t.tup_inserted,
			t.tup_updated,
			t.tup_deleted
		FROM db_stats db
		JOIN cache_stats c ON db.datname = c.datname
		JOIN txn_stats t ON db.datname = t.datname
		LIMIT $1
	`
	
	rows, err := sourceDB.QueryContext(ctx, query, c.config.Monitoring.MaxDatabases)
	if err != nil {
		return err
	}
	defer rows.Close()

	var resources []models.ResourceUsage
	for rows.Next() {
		var r models.ResourceUsage
		r.CheckTime = time.Now()
		
		var (
			totalTransactions int64
			tupFetched       int64
			tupInserted      int64
			tupUpdated       int64
			tupDeleted       int64
		)
		
		if err := rows.Scan(
			&r.DBName,
			&r.ActiveConnections,
			&r.CacheHitRatio,
			&totalTransactions,
			&tupFetched,
			&tupInserted,
			&tupUpdated,
			&tupDeleted,
		); err != nil {
			c.logger.Warnf("Failed to scan resource stats: %v", err)
			continue
		}
		
		r.MaxConnections = maxConnections
		r.ConnectionUsagePct = float64(r.ActiveConnections) / float64(maxConnections) * 100
		r.TransactionsPerSec = float64(totalTransactions) / uptime
		r.TuplesFetchedPerSec = float64(tupFetched) / uptime
		r.TuplesInsertedPerSec = float64(tupInserted) / uptime
		r.TuplesUpdatedPerSec = float64(tupUpdated) / uptime
		r.TuplesDeletedPerSec = float64(tupDeleted) / uptime
		
		resources = append(resources, r)
	}

	if err := c.repo.SaveResourceUsage(ctx, resources); err != nil {
		return err
	}

	c.logger.Infof("Collected resource usage for %d databases", len(resources))
	return nil
}