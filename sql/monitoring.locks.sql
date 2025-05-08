CREATE TABLE monitoring.locks (
    id SERIAL PRIMARY KEY,
    check_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    db_name VARCHAR(255) NOT NULL,
    blocked_pid INTEGER NOT NULL,
    blocked_user VARCHAR(255) NOT NULL,
    blocked_query TEXT,
    blocking_pid INTEGER NOT NULL,
    blocking_user VARCHAR(255) NOT NULL,
    blocking_query TEXT,
    lock_type VARCHAR(255) NOT NULL,
    mode VARCHAR(255) NOT NULL,
    duration INTERVAL NOT NULL
);

CREATE INDEX idx_locks_time ON monitoring.locks(check_time);
CREATE INDEX idx_locks_dbname ON monitoring.locks(db_name);
CREATE INDEX idx_locks_blocked ON monitoring.locks(blocked_pid);
CREATE INDEX idx_locks_blocking ON monitoring.locks(blocking_pid);

SELECT
    blocked.datname AS db_name,
    blocked_locks.pid AS blocked_pid,
    blocked_activity.usename AS blocked_user,
    blocked_activity.query AS blocked_query,
    blocking_locks.pid AS blocking_pid,
    blocking_activity.usename AS blocking_user,
    blocking_activity.query AS blocking_query,
    blocked_locks.locktype AS lock_type,
    blocked_locks.mode AS mode,
    now() - blocked_activity.query_start AS duration
FROM
    pg_catalog.pg_locks blocked_locks
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
WHERE
    NOT blocked_locks.granted;