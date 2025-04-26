package db

import (
    "fmt"
    "os"

    "github.com/IvanSt1/ctfad/otkritki/backend/core/models"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var db *gorm.DB

// init устанавливает соединение с БД и выполняет авто-миграцию
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
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect to database: " + err.Error())
    }

    // Миграция моделей GiftCard и User
    if err := db.AutoMigrate(&models.GiftCard{}, &models.User{}); err != nil {
        panic("failed to auto-migrate models: " + err.Error())
    }
}

// GetDB возвращает единый экземпляр *gorm.DB
func GetDB() *gorm.DB {
    return db
}
