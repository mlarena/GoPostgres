CREATE TABLE monitoring.table_stats (
    id SERIAL PRIMARY KEY,
    check_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    db_name VARCHAR(255) NOT NULL,
    schema_name VARCHAR(255) NOT NULL,
    table_name VARCHAR(255) NOT NULL,
    total_size_pretty TEXT NOT NULL,
    total_size_bytes BIGINT NOT NULL,
    row_count INTEGER NOT NULL
);

CREATE INDEX idx_table_stats_time ON monitoring.table_stats(check_time);
CREATE INDEX idx_table_stats_dbname ON monitoring.table_stats(db_name);
CREATE INDEX idx_table_stats_schema ON monitoring.table_stats(schema_name);
CREATE INDEX idx_table_stats_table ON monitoring.table_stats(table_name);

   -----------------

SELECT
    current_database() AS db_name,
    nspname AS schema_name,
    relname AS table_name,
    pg_size_pretty(pg_total_relation_size(C.oid)) AS total_size_pretty,
    pg_total_relation_size(C.oid) AS total_size_bytes,
    (xpath('/row/cnt/text()', query_to_xml(format('SELECT COUNT(*) AS cnt FROM %I.%I', nspname, relname), false, true, '')))[1]::text::int AS row_count
FROM
    pg_catalog.pg_class C
    JOIN pg_catalog.pg_namespace N ON N.oid = C.relnamespace
WHERE
    C.relkind = 'r'
    AND nspname NOT IN ('pg_catalog', 'information_schema')
    AND nspname !~ '^pg_toast'
ORDER BY
    pg_total_relation_size(C.oid) DESC;

    -----------------


 INSERT INTO monitoring.table_stats (
    check_time, db_name, schema_name, table_name, 
    total_size_pretty, total_size_bytes, row_count
)
SELECT
    NOW(),
    current_database(),
    nspname,
    relname,
    pg_size_pretty(pg_total_relation_size(C.oid)),
    pg_total_relation_size(C.oid),
    (xpath('/row/cnt/text()', query_to_xml(format('SELECT COUNT(*) AS cnt FROM %I.%I', nspname, relname), false, true, '')))[1]::text::int
FROM
    pg_catalog.pg_class C
    JOIN pg_catalog.pg_namespace N ON N.oid = C.relnamespace
WHERE
    C.relkind = 'r'
    AND nspname NOT IN ('pg_catalog', 'information_schema')
    AND nspname !~ '^pg_toast';   