CREATE TABLE monitoring.index_stats (
    id SERIAL PRIMARY KEY,
    check_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    db_name VARCHAR(255) NOT NULL,
    schema_name VARCHAR(255) NOT NULL,
    table_name VARCHAR(255) NOT NULL,
    index_name VARCHAR(255) NOT NULL,
    index_size_pretty TEXT NOT NULL,
    index_size_bytes BIGINT NOT NULL,
    index_scans BIGINT NOT NULL,
    index_tup_read BIGINT NOT NULL,
    index_tup_fetch BIGINT NOT NULL
);

CREATE INDEX idx_index_stats_time ON monitoring.index_stats(check_time);
CREATE INDEX idx_index_stats_dbname ON monitoring.index_stats(db_name);
CREATE INDEX idx_index_stats_table ON monitoring.index_stats(table_name);
CREATE INDEX idx_index_stats_index ON monitoring.index_stats(index_name);


SELECT
    current_database() AS db_name,
    n.nspname AS schema_name,
    t.relname AS table_name,
    i.relname AS index_name,
    pg_size_pretty(pg_relation_size(i.oid)) AS index_size_pretty,
    pg_relation_size(i.oid) AS index_size_bytes,
    ix.idx_scan AS index_scans,
    ix.idx_tup_read AS index_tup_read,
    ix.idx_tup_fetch AS index_tup_fetch
FROM
    pg_class t
    JOIN pg_index x ON t.oid = x.indrelid
    JOIN pg_class i ON i.oid = x.indexrelid
    JOIN pg_namespace n ON n.oid = t.relnamespace
    JOIN pg_stat_all_indexes ix ON ix.indexrelid = i.oid
WHERE
    t.relkind = 'r'
    AND n.nspname NOT IN ('pg_catalog', 'information_schema')
ORDER BY
    pg_relation_size(i.oid) DESC;