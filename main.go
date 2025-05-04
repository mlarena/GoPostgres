package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
)

func listDatabases() {
    // Параметры подключения к PostgreSQL
    host := "localhost"
    port := "5432"
    user := "postgres"
    password := "12345678"
    dbname := "postgres" // Имя базы данных для подключения (обычно "postgres")

    // Строка подключения
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    // Подключение к базе данных
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error opening database: %q", err)
    }
    defer db.Close()

    // Проверка подключения
    err = db.Ping()
    if err != nil {
        log.Fatalf("Error connecting to the database: %q", err)
    }
    fmt.Println("Successfully connected to the database!")

    // Запрос списка баз данных
    rows, err := db.Query("SELECT datname FROM pg_database WHERE datistemplate = false;")
    if err != nil {
        log.Fatalf("Error querying databases: %q", err)
    }
    defer rows.Close()

    // Вывод списка баз данных
    fmt.Println("List of databases:")
    for rows.Next() {
        var dbName string
        if err := rows.Scan(&dbName); err != nil {
            log.Fatalf("Error scanning row: %q", err)
        }
        fmt.Println(dbName)
    }

    // Проверка ошибок после завершения итерации
    if err = rows.Err(); err != nil {
        log.Fatalf("Error during row iteration: %q", err)
    }
}

func getPostgresVersion() {
    // Параметры подключения к PostgreSQL
    host := "localhost"
    port := "5432"
    user := "postgres"
    password := "12345678"
    dbname := "postgres" // Имя базы данных для подключения (обычно "postgres")

    // Строка подключения
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    // Подключение к базе данных
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error opening database: %q", err)
    }
    defer db.Close()

    // Проверка подключения
    err = db.Ping()
    if err != nil {
        log.Fatalf("Error connecting to the database: %q", err)
    }
    fmt.Println("Successfully connected to the database!")

    // Запрос версии PostgreSQL
    var version string
    err = db.QueryRow("SELECT version();").Scan(&version)
    if err != nil {
        log.Fatalf("Error querying PostgreSQL version: %q", err)
    }

    // Вывод версии PostgreSQL
    fmt.Println("PostgreSQL version:", version)
}

func getDatabaseSizes() {
    // Параметры подключения к PostgreSQL
    host := "localhost"
    port := "5432"
    user := "postgres"
    password := "12345678"
    dbname := "postgres" // Имя базы данных для подключения (обычно "postgres")

    // Строка подключения
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    // Подключение к базе данных
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error opening database: %q", err)
    }
    defer db.Close()

    // Проверка подключения
    err = db.Ping()
    if err != nil {
        log.Fatalf("Error connecting to the database: %q", err)
    }
    fmt.Println("Successfully connected to the database!")

    // Запрос списка баз данных и их размеров
    rows, err := db.Query(`
        SELECT datname AS db_name, pg_size_pretty(pg_database_size(datname)) AS db_size
        FROM pg_database
        ORDER BY pg_database_size(datname) DESC;
    `)
    if err != nil {
        log.Fatalf("Error querying database sizes: %q", err)
    }
    defer rows.Close()

    // Вывод списка баз данных и их размеров
    fmt.Println("List of databases and their sizes:")
    for rows.Next() {
        var dbName, dbSize string
        if err := rows.Scan(&dbName, &dbSize); err != nil {
            log.Fatalf("Error scanning row: %q", err)
        }
        fmt.Printf("%s: %s\n", dbName, dbSize)
    }

    // Проверка ошибок после завершения итерации
    if err = rows.Err(); err != nil {
        log.Fatalf("Error during row iteration: %q", err)
    }
}

func getTableInfo(dbname string) {
    // Параметры подключения к PostgreSQL
    host := "localhost"
    port := "5432"
    user := "postgres"
    password := "12345678"

    // Строка подключения
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    // Подключение к базе данных
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error opening database: %q", err)
    }
    defer db.Close()

    // Проверка подключения
    err = db.Ping()
    if err != nil {
        log.Fatalf("Error connecting to the database: %q", err)
    }
    fmt.Println("Successfully connected to the database!")

    // Запрос информации о таблицах
    rows, err := db.Query(`
        SELECT
            t.table_name AS "Таблица",
            pg_size_pretty(pg_total_relation_size(quote_ident(t.table_schema) || '.' || quote_ident(t.table_name))) AS "Размер",
            (xpath('/row/cnt/text()', query_to_xml(format('SELECT COUNT(*) AS cnt FROM %I.%I', t.table_schema, t.table_name), false, true, '')))[1]::text::int AS "Количество строк"
        FROM
            information_schema.tables t
        WHERE
            t.table_schema NOT IN ('pg_catalog', 'information_schema')
            AND t.table_type = 'BASE TABLE'
        ORDER BY
            pg_total_relation_size(quote_ident(t.table_schema) || '.' || quote_ident(t.table_name)) DESC;
    `)
    if err != nil {
        log.Fatalf("Error querying table info: %q", err)
    }
    defer rows.Close()

    // Вывод информации о таблицах
    fmt.Println("List of tables and their sizes:")
    for rows.Next() {
        var tableName, tableSize string
        var rowCount int
        if err := rows.Scan(&tableName, &tableSize, &rowCount); err != nil {
            log.Fatalf("Error scanning row: %q", err)
        }
        fmt.Printf("%s: %s, %d rows\n", tableName, tableSize, rowCount)
    }

    // Проверка ошибок после завершения итерации
    if err = rows.Err(); err != nil {
        log.Fatalf("Error during row iteration: %q", err)
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <command> [dbname]")
        fmt.Println("Available commands:")
        fmt.Println("  dblist   - List all databases")
        fmt.Println("  version  - Get PostgreSQL version")
        fmt.Println("  dbsizes  - List databases and their sizes")
        fmt.Println("  tblinfo  - Get table info for a specific database")
        return
    }

    command := os.Args[1]
    switch command {
    case "dblist":
        listDatabases()
    case "version":
        getPostgresVersion()
    case "dbsizes":
        getDatabaseSizes()
    case "tblinfo":
        if len(os.Args) < 3 {
            fmt.Println("Usage: go run main.go tblinfo <dbname>")
            return
        }
        dbname := os.Args[2]
        getTableInfo(dbname)
    default:
        fmt.Printf("Unknown command: %s\n", command)
        fmt.Println("Available commands:")
        fmt.Println("  dblist   - List all databases")
        fmt.Println("  version  - Get PostgreSQL version")
        fmt.Println("  dbsizes  - List databases and their sizes")
        fmt.Println("  tblinfo  - Get table info for a specific database")
    }
}
