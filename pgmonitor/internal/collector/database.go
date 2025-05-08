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

type DatabaseCollector struct {
	repo   *repository.Repository
	logger *logger.Logger
	config *config.Config
}

func NewDatabaseCollector(repo *repository.Repository, logger *logger.Logger, cfg *config.Config) *DatabaseCollector {
	return &DatabaseCollector{
		repo:   repo,
		logger: logger,
		config: cfg,
	}
}

func (c *DatabaseCollector) Collect(ctx context.Context, sourceDB *sql.DB) error {
	c.logger.Info("Collecting database statistics")
	
	query := `
		SELECT 
			datname, 
			pg_size_pretty(pg_database_size(datname)), 
			pg_database_size(datname), 
			datcollate, 
			datconnlimit, 
			datallowconn
		FROM pg_database
		LIMIT $1
	`
	
	rows, err := sourceDB.QueryContext(ctx, query, c.config.Monitoring.MaxDatabases)
	if err != nil {
		c.logger.Errorf("Failed to query database stats: %v", err)
		return err
	}
	defer rows.Close()
	
	var stats []models.DatabaseStat
	for rows.Next() {
		var stat models.DatabaseStat
		stat.CheckTime = time.Now()
		
		if err := rows.Scan(
			&stat.DBName,
			&stat.SizePretty,
			&stat.SizeBytes,
			&stat.Collation,
			&stat.ConnectionLimit,
			&stat.ConnectionsAllowed,
		); err != nil {
			c.logger.Warnf("Failed to scan database stat: %v", err)
			continue
		}
		
		stats = append(stats, stat)
	}
	
	if err := rows.Err(); err != nil {
		c.logger.Errorf("Error iterating database stats: %v", err)
		return err
	}
	
	if err := c.repo.SaveDatabaseStats(ctx, stats); err != nil {
		c.logger.Errorf("Failed to save database stats: %v", err)
		return err
	}
	
	c.logger.Infof("Successfully collected stats for %d databases", len(stats))
	return nil
}