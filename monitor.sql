DO $$
DECLARE
    db_record RECORD;
    query_text TEXT;
    result RECORD;
    exec_id BIGINT;
    start_time TIMESTAMP;
    total_tables INTEGER := 0;
BEGIN
    -- Фиксируем начало выполнения
    INSERT INTO workdb.monitoring.execution_history (
        execution_time, duration, databases_scanned, tables_scanned, success, error_message
    ) VALUES (
        NOW(), NULL, 0, 0, TRUE, NULL
    ) RETURNING id INTO exec_id;
    
    start_time := NOW();
    
    -- Собираем глобальную статистику из postgres
    -- database_stats
    INSERT INTO workdb.monitoring.database_stats
    SELECT NOW(), datname, pg_size_pretty(pg_database_size(datname)), 
           pg_database_size(datname), datcollate, datconnlimit, datallowconn
    FROM pg_database;
    
    -- long_running_queries
    INSERT INTO workdb.monitoring.long_running_queries
    SELECT NOW(), datname, pid, usename, application_name, client_addr, 
           backend_start, query_start, NOW() - query_start, query, state
    FROM pg_stat_activity
    WHERE state = 'active' AND query_start IS NOT NULL AND NOW() - query_start > interval '5 minutes';
    
    -- locks
    INSERT INTO workdb.monitoring.locks
    SELECT NOW(), blocked.datname, blocked_locks.pid, blocked_activity.usename, 
           blocked_activity.query, blocking_locks.pid, blocking_activity.usename,
           blocking_activity.query, blocked_locks.locktype, blocked_locks.mode, 
           NOW() - blocked_activity.query_start
    FROM pg_catalog.pg_locks blocked_locks
    JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
    JOIN pg_catalog.pg_locks blocking_locks ON blocking_locks.locktype = blocked_locks.locktype
        AND blocking_locks.DATABASE IS NOT DISTINCT FROM blocked_locks.DATABASE
        AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
        AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
        AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
        AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
        AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
        AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
        AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
        AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
        AND blocking_locks.pid != blocked_locks.pid
    JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
    WHERE NOT blocked_locks.granted;
    
    -- resource_usage
    WITH db_stats AS (
        SELECT datname, COUNT(*) AS active_connections FROM pg_stat_activity GROUP BY datname
    ),
    cache_stats AS (
        SELECT datname, sum(blks_hit)/nullif(sum(blks_hit+blks_read),0)*100 AS cache_hit_ratio 
        FROM pg_stat_database GROUP BY datname
    ),
    txn_stats AS (
        SELECT datname, xact_commit+xact_rollback AS total_transactions,
               tup_fetched, tup_inserted, tup_updated, tup_deleted
        FROM pg_stat_database
    )
    INSERT INTO workdb.monitoring.resource_usage
    SELECT NOW(), db.datname, db.active_connections, 
           (SELECT setting FROM pg_settings WHERE name = 'max_connections')::int,
           db.active_connections*100.0/(SELECT setting FROM pg_settings WHERE name = 'max_connections')::int,
           c.cache_hit_ratio,
           t.total_transactions/extract(epoch FROM (NOW()-pg_postmaster_start_time())),
           t.tup_fetched/extract(epoch FROM (NOW()-pg_postmaster_start_time())),
           t.tup_inserted/extract(epoch FROM (NOW()-pg_postmaster_start_time())),
           t.tup_updated/extract(epoch FROM (NOW()-pg_postmaster_start_time())),
           t.tup_deleted/extract(epoch FROM (NOW()-pg_postmaster_start_time()))
    FROM db_stats db
    JOIN cache_stats c ON db.datname = c.datname
    JOIN txn_stats t ON db.datname = t.datname;
    
    -- Обходим все базы данных для сбора специфической статистики
    FOR db_record IN SELECT datname FROM pg_database WHERE datallowconn AND datname NOT IN ('template0', 'template1', 'postgres', 'workdb')
    LOOP
        -- table_stats
        query_text := format('
            INSERT INTO workdb.monitoring.table_stats
            SELECT 
                NOW(), 
                %L, 
                nspname, 
                relname, 
                pg_size_pretty(pg_total_relation_size(C.oid)), 
                pg_total_relation_size(C.oid),
                (xpath(''/row/cnt/text()'', query_to_xml(format(''SELECT COUNT(*) AS cnt FROM %%I.%%I'', nspname, relname), false, true, '''')))[1]::text::int
            FROM pg_catalog.pg_class C
            JOIN pg_catalog.pg_namespace N ON N.oid = C.relnamespace
            WHERE C.relkind = ''r'' AND nspname NOT IN (''pg_catalog'', ''information_schema'') AND nspname !~ ''^pg_toast''
        ', db_record.datname);
        
        PERFORM dblink_exec('dbname=' || db_record.datname || ' user=monitoring_user password=monitoring_password', query_text);
        GET DIAGNOSTICS total_tables = ROW_COUNT;
        
        -- index_stats
        query_text := format('
            INSERT INTO workdb.monitoring.index_stats
            SELECT 
                NOW(), 
                %L, 
                n.nspname, 
                t.relname, 
                i.relname, 
                pg_size_pretty(pg_relation_size(i.oid)), 
                pg_relation_size(i.oid), 
                ix.idx_scan, 
                ix.idx_tup_read, 
                ix.idx_tup_fetch
            FROM pg_class t
            JOIN pg_index x ON t.oid = x.indrelid
            JOIN pg_class i ON i.oid = x.indexrelid
            JOIN pg_namespace n ON n.oid = t.relnamespace
            JOIN pg_stat_all_indexes ix ON ix.indexrelid = i.oid
            WHERE t.relkind = ''r'' AND n.nspname NOT IN (''pg_catalog'', ''information_schema'')
        ', db_record.datname);
        
        PERFORM dblink_exec('dbname=' || db_record.datname || ' user=monitoring_user password=monitoring_password', query_text);
    END LOOP;
    
    -- Обновляем запись о выполнении
    UPDATE workdb.monitoring.execution_history 
    SET duration = NOW() - start_time,
        databases_scanned = (SELECT COUNT(*) FROM pg_database WHERE datallowconn AND datname NOT IN ('template0', 'template1')),
        tables_scanned = total_tables
    WHERE id = exec_id;
EXCEPTION WHEN OTHERS THEN
    UPDATE workdb.monitoring.execution_history 
    SET success = FALSE,
        error_message = SQLERRM
    WHERE id = exec_id;
END $$;