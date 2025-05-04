package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
)

func main() {
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
