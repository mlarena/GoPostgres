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

type QueryCollector struct {
	repo   *repository.Repository
	logger *logger.Logger
	config *config.Config
}

func NewQueryCollector(repo *repository.Repository, logger *logger.Logger, cfg *config.Config) *QueryCollector {
	return &QueryCollector{
		repo:   repo,
		logger: logger,
		config: cfg,
	}
}

func (c *QueryCollector) Collect(ctx context.Context, sourceDB *sql.DB) error {
	c.logger.Info("Collecting long running queries")
	
	threshold, err := c.config.Monitoring.GetLongQueryThreshold()
	if err != nil {
		return err
	}

	query := `
		SELECT 
			datname, pid, usename, application_name,
			client_addr, backend_start, query_start,
			query, state
		FROM pg_stat_activity
		WHERE state = 'active' 
		AND query_start IS NOT NULL 
		AND now() - query_start > $1
	`
	
	rows, err := sourceDB.QueryContext(ctx, query, threshold)
	if err != nil {
		return err
	}
	defer rows.Close()

	var queries []models.LongRunningQuery
	for rows.Next() {
		var q models.LongRunningQuery
		q.CheckTime = time.Now()
		
		if err := rows.Scan(
			&q.DBName,
			&q.PID,
			&q.Username,
			&q.Application,
			&q.ClientAddr,
			&q.BackendStart,
			&q.QueryStart,
			&q.Query,
			&q.State,
		); err != nil {
			c.logger.Warnf("Failed to scan query: %v", err)
			continue
		}
		
		q.Duration = time.Since(q.QueryStart)
		queries = append(queries, q)
	}

	if err := c.repo.SaveLongRunningQueries(ctx, queries); err != nil {
		return err
	}

	c.logger.Infof("Collected %d long running queries", len(queries))
	return nil
}