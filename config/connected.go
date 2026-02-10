package config

import (
    "database/sql"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL is not set")
    }

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Failed to open database:", err)
    }

    if err := db.Ping(); err != nil {
        log.Fatal("Database unreachable:", err)
    }

    log.Println("Database connected successfully")
    return db
}