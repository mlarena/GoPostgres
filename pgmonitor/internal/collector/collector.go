package collector

import (
	"context"
	"database/sql"
	
	"pgmonitor/internal/repository"
	"pgmonitor/pkg/logger"
	"pgmonitor/config"
)

type Collector struct {
	dbCollector    *DatabaseCollector
	queryCollector *QueryCollector
	lockCollector  *LockCollector
	resCollector   *ResourceCollector
	execCollector  *ExecutionCollector
	logger         *logger.Logger
	config         *config.Config
}

func NewCollector(repo *repository.Repository, logger *logger.Logger, cfg *config.Config) *Collector {
	return &Collector{
		dbCollector:    NewDatabaseCollector(repo, logger, cfg),
		queryCollector: NewQueryCollector(repo, logger, cfg),
		lockCollector:  NewLockCollector(repo, logger, cfg),
		resCollector:   NewResourceCollector(repo, logger, cfg),
		execCollector:  NewExecutionCollector(repo, logger),
		logger:         logger,
		config:         cfg,
	}
}

func (c *Collector) Run(ctx context.Context) error {
	// Начинаем выполнение
	execID, err := c.execCollector.StartExecution(ctx)
	if err != nil {
		return err
	}
	
	var (
		dbsScanned    int
		tablesScanned int
		success       = true
		errMsg       string
	)
	
	defer func() {
		// Завершаем выполнение
		if err := c.execCollector.EndExecution(ctx, execID, dbsScanned, tablesScanned, success, errMsg); err != nil {
			c.logger.Errorf("Failed to finalize execution: %v", err)
		}
	}()
	
	// Подключаемся к source DB
	sourceDB, err := sql.Open("postgres", c.config.Database["source"].ConnectionString())
	if err != nil {
		errMsg = err.Error()
		success = false
		return err
	}
	defer sourceDB.Close()
	
	// Собираем данные
	if err := c.dbCollector.Collect(ctx, sourceDB); err != nil {
		errMsg = err.Error()
		success = false
	}
	
	if err := c.queryCollector.Collect(ctx, sourceDB); err != nil {
		errMsg = err.Error()
		success = false
	}
	
	if err := c.lockCollector.Collect(ctx, sourceDB); err != nil {
		errMsg = err.Error()
		success = false
	}
	
	if err := c.resCollector.Collect(ctx, sourceDB); err != nil {
		errMsg = err.Error()
		success = false
	}
	
	return nil
}