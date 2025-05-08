CREATE TABLE monitoring.resource_usage (
    id SERIAL PRIMARY KEY,
    check_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    db_name VARCHAR(255) NOT NULL,
    active_connections INTEGER NOT NULL,
    max_connections INTEGER NOT NULL,
    connection_usage_pct NUMERIC(5,2) NOT NULL,
    cache_hit_ratio NUMERIC(5,2) NOT NULL,
    transactions_per_sec NUMERIC(10,2) NOT NULL,
    tuples_fetched_per_sec NUMERIC(10,2) NOT NULL,
    tuples_inserted_per_sec NUMERIC(10,2) NOT NULL,
    tuples_updated_per_sec NUMERIC(10,2) NOT NULL,
    tuples_deleted_per_sec NUMERIC(10,2) NOT NULL
);

CREATE INDEX idx_resource_usage_time ON monitoring.resource_usage(check_time);
CREATE INDEX idx_resource_usage_dbname ON monitoring.resource_usage(db_name);