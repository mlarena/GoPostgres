package collector

import (
	"context"
	"time"
	
	"pgmonitor/internal/repository"
	"pgmonitor/pkg/logger"
)

type ExecutionCollector struct {
	repo   *repository.Repository
	logger *logger.Logger
}

func NewExecutionCollector(repo *repository.Repository, logger *logger.Logger) *ExecutionCollector {
	return &ExecutionCollector{
		repo:   repo,
		logger: logger,
	}
}

func (c *ExecutionCollector) StartExecution(ctx context.Context) (int64, error) {
	execID, err := c.repo.StartExecution(ctx)
	if err != nil {
		c.logger.Errorf("Failed to start execution: %v", err)
		return 0, err
	}
	
	c.logger.Infof("Started new monitoring execution with ID: %d", execID)
	return execID, nil
}

func (c *ExecutionCollector) EndExecution(ctx context.Context, execID int64, dbsScanned, tablesScanned int, success bool, errMsg string) error {
	duration := time.Since(time.Now())
	
	err := c.repo.EndExecution(ctx, execID, duration, dbsScanned, tablesScanned, success, errMsg)
	if err != nil {
		c.logger.Errorf("Failed to end execution %d: %v", execID, err)
		return err
	}
	
	if success {
		c.logger.Infof("Completed execution %d successfully. Duration: %v", execID, duration)
	} else {
		c.logger.Warnf("Completed execution %d with errors. Duration: %v, Error: %s", execID, duration, errMsg)
	}
	
	return nil
}