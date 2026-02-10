package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

func LoadEnv() {
    if err := godotenv.Load(); 
	err != nil {
        log.Println(".env tidak ditemukan, menggunakan environment bawaan")
    }
}

func GetEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}
