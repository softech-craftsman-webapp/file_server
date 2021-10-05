package config

import (
	model "file_server/model"
)

func Migrate() {
	db := GetDB()

	// Auto Migration
	db.AutoMigrate(&model.File{})

	CloseDB(db).Close()
}
