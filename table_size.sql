SELECT
    nspname AS schema_name,
    relname AS table_name,
    pg_size_pretty(pg_total_relation_size(C.oid)) AS total_size,
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