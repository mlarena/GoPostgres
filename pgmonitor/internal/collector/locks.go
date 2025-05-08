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

type LockCollector struct {
	repo   *repository.Repository
	logger *logger.Logger
	config *config.Config
}

func NewLockCollector(repo *repository.Repository, logger *logger.Logger, cfg *config.Config) *LockCollector {
	return &LockCollector{
		repo:   repo,
		logger: logger,
		config: cfg,
	}
}

func (c *LockCollector) Collect(ctx context.Context, sourceDB *sql.DB) error {
	c.logger.Info("Collecting lock information")
	
	query := `
		SELECT
			blocked.datname,
			blocked_locks.pid AS blocked_pid,
			blocked_activity.usename AS blocked_user,
			blocked_activity.query AS blocked_query,
			blocking_locks.pid AS blocking_pid,
			blocking_activity.usename AS blocking_user,
			blocking_activity.query AS blocking_query,
			blocked_locks.locktype AS lock_type,
			blocked_locks.mode AS mode
		FROM pg_catalog.pg_locks blocked_locks
		JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
		JOIN pg_catalog.pg_locks blocking_locks ON blocking_locks.locktype = blocked_locks.locktype
			AND blocking_locks.DATABASE IS NOT DISTINCT FROM blocked_locks.DATABASE
			AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
			AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
			AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
			AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
			AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
			AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
			AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
			AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
			AND blocking_locks.pid != blocked_locks.pid
		JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
		WHERE NOT blocked_locks.granted
	`
	
	rows, err := sourceDB.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var locks []models.Lock
	for rows.Next() {
		var l models.Lock
		l.CheckTime = time.Now()
		
		if err := rows.Scan(
			&l.DBName,
			&l.BlockedPID,
			&l.BlockedUser,
			&l.BlockedQuery,
			&l.BlockingPID,
			&l.BlockingUser,
			&l.BlockingQuery,
			&l.LockType,
			&l.Mode,
		); err != nil {
			c.logger.Warnf("Failed to scan lock: %v", err)
			continue
		}
		
		// Получаем продолжительность блокировки
		var duration time.Duration
		err = sourceDB.QueryRowContext(ctx, 
			"SELECT now() - query_start FROM pg_stat_activity WHERE pid = $1", 
			l.BlockedPID).Scan(&duration)
		if err != nil {
			duration = 0
		}
		l.Duration = duration
		
		locks = append(locks, l)
	}

	if err := c.repo.SaveLocks(ctx, locks); err != nil {
		return err
	}

	c.logger.Infof("Collected %d lock events", len(locks))
	return nil
}