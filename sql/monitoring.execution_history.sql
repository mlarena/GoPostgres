CREATE TABLE monitoring.execution_history (
    id SERIAL PRIMARY KEY,
    execution_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    duration INTERVAL NOT NULL,
    databases_scanned INTEGER NOT NULL,
    tables_scanned INTEGER NOT NULL,
    success BOOLEAN NOT NULL,
    error_message TEXT
);

CREATE INDEX idx_execution_history_time ON monitoring.execution_history(execution_time);