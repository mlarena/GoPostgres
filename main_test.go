package main

import (
    "database/sql"
    "fmt"
    "testing"

    _ "github.com/lib/pq"
)

// Юнит-тест для проверки подключения к базе данных
func TestDatabaseConnection(t *testing.T) {
    // Параметры подключения к PostgreSQL
    host := "localhost"
    port := "5432"
    user := "postgres"
    password := "12345678"
    dbname := "postgres" // Имя базы данных для подключения (обычно "postgres")

    // Строка подключения
    connStr := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"

    // Подключение к базе данных
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        t.Fatalf("Error opening database: %v", err)
    }
    defer db.Close()

    // Проверка подключения
    err = db.Ping()
    if err != nil {
        t.Fatalf("Error connecting to the database: %v", err)
    }
}

// Пример использования функции для подключения к базе данных
func ExampleDatabaseConnection() {
    // Параметры подключения к PostgreSQL
    host := "localhost"
    port := "5432"
    user := "postgres"
    password := "12345678"
    dbname := "postgres"// Имя базы данных для подключения (обычно "postgres")

    // Строка подключения
    connStr := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"

    // Подключение к базе данных
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic("Error opening database: " + err.Error())
    }
    defer db.Close()

    // Проверка подключения
    err = db.Ping()
    if err != nil {
        panic("Error connecting to the database: " + err.Error())
    }

    // Пример вывода
    fmt.Println("Successfully connected to the database!")

    // Output:
    // Successfully connected to the database!
}
