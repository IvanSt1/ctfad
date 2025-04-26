package db

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func NewDb(connectionString string, Structs ...any) *DB {
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(Structs...); err != nil {
		panic(err)
	}
	return &DB{
		db,
	}
}

func (db *DB) AddTable(Struct any) (err error) {

	if db.db.Migrator().HasTable(Struct) {
		errors.New("Table for this struct already exists")
	}
	err = db.db.AutoMigrate(Struct)
	return err
}
