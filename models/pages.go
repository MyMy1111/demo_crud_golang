package models

import "gorm.io/gorm"

type Pages struct {
	Id				 uint					`gorm:"primary key;autoincrement" json: "id"`
	Body      *string				`json: "body"`
	Page_id   *string				`json: "page_id"`
}

type UpdatePages struct {
	Body      *string				`json: "body"`
	Page_id   *string				`json: "page_id"`
}

func MigratePages(db *gorm.DB) error {
	err := db.AutoMigrate(&Pages{})
	return err
}