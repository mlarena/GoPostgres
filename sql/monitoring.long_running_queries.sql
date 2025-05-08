CREATE TABLE monitoring.long_running_queries (
    id SERIAL PRIMARY KEY,
    check_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    db_name VARCHAR(255) NOT NULL,
    pid INTEGER NOT NULL,
    usename VARCHAR(255) NOT NULL,
    application_name VARCHAR(255),
    client_addr INET,
    backend_start TIMESTAMP WITH TIME ZONE,
    query_start TIMESTAMP WITH TIME ZONE,
    duration INTERVAL NOT NULL,
    query TEXT NOT NULL,
    state VARCHAR(30) NOT NULL
);

CREATE INDEX idx_long_queries_time ON monitoring.long_running_queries(check_time);
CREATE INDEX idx_long_queries_dbname ON monitoring.long_running_queries(db_name);
CREATE INDEX idx_long_queries_duration ON monitoring.long_running_queries(duration);

SELECT
    datname AS db_name,
    pid,
    usename,
    application_name,
    client_addr,
    backend_start,
    query_start,
    now() - query_start AS duration,
    query,
    state
FROM
    pg_stat_activity
WHERE
    state = 'active'
    AND query_start IS NOT NULL
    AND now() - query_start > interval '5 minutes'
ORDER BY
    duration DESC;