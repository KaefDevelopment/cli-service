package model

import (
	"gorm.io/gorm"
	"log"
)

const (
	createTableQuery = `CREATE TABLE IF NOT EXISTS events (
    	id varchar(255) not null,
    	createdAt varchar(255) not null,
    	type varchar(255) not null,
    	project varchar(255),
    	projectBaseDir varchar(255),
    	language varchar(255),
    	target varchar(255),
    	branch varchat(255),
    	timeZone varchar(255),
    	params text,
    	authKey varchar(255) not null,
    	send bool
);`
)

func CreateTable(db *gorm.DB) error {
	tx := db.Exec(createTableQuery)
	if tx.Error != nil {
		log.Println("create table failed:", tx.Error)
		return tx.Error
	}

	return nil
}
