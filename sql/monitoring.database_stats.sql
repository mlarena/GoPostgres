CREATE TABLE monitoring.database_stats (
    id SERIAL PRIMARY KEY,
    check_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    db_name VARCHAR(255) NOT NULL,
    size_pretty TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,
    db_collation VARCHAR(255) NOT NULL,
    connection_limit INTEGER NOT NULL,
    connections_allowed BOOLEAN NOT NULL  -- Исправленная опечатка
);

CREATE INDEX idx_database_stats_time ON monitoring.database_stats(check_time);
CREATE INDEX idx_database_stats_dbname ON monitoring.database_stats(db_name);


SELECT 
    datname AS db_name,
    pg_size_pretty(pg_database_size(datname)) AS size_pretty,
    pg_database_size(datname) AS size_bytes,
    datcollate AS db_collation,
    datconnlimit AS connection_limit,
    datallowconn AS connections_allowed
FROM 
    pg_database
ORDER BY 
    pg_database_size(datname) DESC;