package db

import (
    "fmt"
    "os"
    "time"

    "github.com/IvanSt1/ctfad/otkritki/backend/core/models"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var db *gorm.DB

// init устанавливает соединение с БД, ожидая готовности сервиса
func init() {
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        os.Getenv("MYSQL_USER"),
        os.Getenv("MYSQL_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )

    var err error
    // Пытаемся подключиться с ретраями
    for i := 0; i < 10; i++ {
        db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
            Logger: logger.Default.LogMode(logger.Silent),
        })
        if err == nil {
            break
        }
        fmt.Printf("[db] connection attempt %d failed: %v\n", i+1, err)
        time.Sleep(5 * time.Second)
    }
    if err != nil {
        panic(fmt.Sprintf("failed to connect to database after retries: %v", err))
    }

    // Миграция моделей
    if err := db.AutoMigrate(&models.GiftCard{}, &models.User{}); err != nil {
        panic(fmt.Sprintf("failed to auto-migrate models: %v", err))
    }
}

// GetDB возвращает единый экземпляр *gorm.DB
func GetDB() *gorm.DB {
    return db
}