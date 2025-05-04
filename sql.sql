SELECT
    relname AS table_name,
    pg_size_pretty(pg_total_relation_size(C.oid)) AS total_size
FROM
    pg_class C
LEFT JOIN
    pg_namespace N ON (N.oid = C.relnamespace)
WHERE
    nspname NOT IN ('pg_catalog', 'information_schema')
    AND C.relkind = 'r'
    AND nspname !~ '^pg_toast'
ORDER BY
    pg_total_relation_size(C.oid) DESC;



SELECT
    table_schema,
    table_name,
    pg_size_pretty(pg_total_relation_size(table_schema || '.' || table_name)) AS total_size
FROM
    information_schema.tables
WHERE
    table_type = 'BASE TABLE'
    AND table_schema NOT IN ('information_schema', 'pg_catalog')
ORDER BY
    pg_total_relation_size(table_schema || '.' || table_name) DESC;

SELECT
    table_schema,
    table_name,
    pg_size_pretty(pg_total_relation_size(table_schema || '.' || table_name)) AS total_size,
    (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = t.table_schema AND table_name = t.table_name) AS row_count
FROM
    information_schema.tables t
WHERE
    table_type = 'BASE TABLE'
    AND table_schema NOT IN ('information_schema', 'pg_catalog')
ORDER BY
    pg_total_relation_size(table_schema || '.' || table_name) DESC;
