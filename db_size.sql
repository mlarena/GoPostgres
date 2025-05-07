SELECT 
    datname AS "Имя базы данных",
    pg_size_pretty(pg_database_size(datname)) AS "Размер",
    pg_database_size(datname) AS "Размер в байтах",
    datcollate AS "Колляция",
    datconnlimit AS "Лимит подключений",
    datallowconn AS "Разрешены подключения"
FROM 
    pg_database
ORDER BY 
    pg_database_size(datname) DESC;