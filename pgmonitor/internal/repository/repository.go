package repository

import (
	"context"
	"database/sql"
	"time"

	"pgmonitor/internal/models"
	"pgmonitor/pkg/logger"
)

type Repository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewRepository(db *sql.DB, logger *logger.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) SaveDatabaseStats(ctx context.Context, stats []models.DatabaseStat) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO monitoring.database_stats (
			check_time, db_name, size_pretty, size_bytes,
			db_collation, connection_limit, connections_allowed
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, stat := range stats {
		if _, err := stmt.ExecContext(
			ctx,
			stat.CheckTime,
			stat.DBName,
			stat.SizePretty,
			stat.SizeBytes,
			stat.Collation,
			stat.ConnectionLimit,
			stat.ConnectionsAllowed,
		); err != nil {
			r.logger.Warnf("Failed to insert database stat for %s: %v", stat.DBName, err)
			continue
		}
	}

	return tx.Commit()
}

func (r *Repository) SaveLongRunningQueries(ctx context.Context, queries []models.LongRunningQuery) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO monitoring.long_running_queries (
			check_time, db_name, pid, usename, application_name,
			client_addr, backend_start, query_start, duration,
			query, state
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, q := range queries {
		if _, err := stmt.ExecContext(
			ctx,
			q.CheckTime,
			q.DBName,
			q.PID,
			q.Username,
			q.Application,
			q.ClientAddr,
			q.BackendStart,
			q.QueryStart,
			q.Duration,
			q.Query,
			q.State,
		); err != nil {
			r.logger.Warnf("Failed to insert long running query for PID %d: %v", q.PID, err)
			continue
		}
	}

	return tx.Commit()
}

func (r *Repository) SaveLocks(ctx context.Context, locks []models.Lock) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO monitoring.locks (
			check_time, db_name, blocked_pid, blocked_user,
			blocked_query, blocking_pid, blocking_user,
			blocking_query, lock_type, mode, duration
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, l := range locks {
		if _, err := stmt.ExecContext(
			ctx,
			l.CheckTime,
			l.DBName,
			l.BlockedPID,
			l.BlockedUser,
			l.BlockedQuery,
			l.BlockingPID,
			l.BlockingUser,
			l.BlockingQuery,
			l.LockType,
			l.Mode,
			l.Duration,
		); err != nil {
			r.logger.Warnf("Failed to insert lock for blocked PID %d: %v", l.BlockedPID, err)
			continue
		}
	}

	return tx.Commit()
}

func (r *Repository) SaveResourceUsage(ctx context.Context, resources []models.ResourceUsage) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO monitoring.resource_usage (
			check_time, db_name, active_connections, max_connections,
			connection_usage_pct, cache_hit_ratio, transactions_per_sec,
			tuples_fetched_per_sec, tuples_inserted_per_sec,
			tuples_updated_per_sec, tuples_deleted_per_sec
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, res := range resources {
		if _, err := stmt.ExecContext(
			ctx,
			res.CheckTime,
			res.DBName,
			res.ActiveConnections,
			res.MaxConnections,
			res.ConnectionUsagePct,
			res.CacheHitRatio,
			res.TransactionsPerSec,
			res.TuplesFetchedPerSec,
			res.TuplesInsertedPerSec,
			res.TuplesUpdatedPerSec,
			res.TuplesDeletedPerSec,
		); err != nil {
			r.logger.Warnf("Failed to insert resource usage for %s: %v", res.DBName, err)
			continue
		}
	}

	return tx.Commit()
}

func (r *Repository) StartExecution(ctx context.Context) (int64, error) {
	var execID int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO monitoring.execution_history (
			execution_time, duration, databases_scanned, 
			tables_scanned, success, error_message
		) VALUES (NOW(), NULL, 0, 0, TRUE, NULL)
		RETURNING id
	`).Scan(&execID)
	
	if err != nil {
		r.logger.Errorf("Failed to start execution: %v", err)
		return 0, err
	}
	
	return execID, nil
}

func (r *Repository) EndExecution(ctx context.Context, execID int64, duration time.Duration, 
	dbsScanned, tablesScanned int, success bool, errMsg string) error {
	
	_, err := r.db.ExecContext(ctx, `
		UPDATE monitoring.execution_history 
		SET duration = $1,
			databases_scanned = $2,
			tables_scanned = $3,
			success = $4,
			error_message = $5
		WHERE id = $6
	`, duration, dbsScanned, tablesScanned, success, errMsg, execID)
	
	if err != nil {
		r.logger.Errorf("Failed to update execution history: %v", err)
		return err
	}
	
	return nil
}