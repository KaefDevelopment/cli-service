package model

import (
	"gorm.io/gorm"
	"log/slog"
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
    	branch varchar(255),
    	timeZone varchar(255),
    	params text,
    	authKey varchar(255) not null,
    	send bool
);`
)

func CreateTable(db *gorm.DB) error {
	tx := db.Exec(createTableQuery)
	if tx.Error != nil {
		slog.Error("create table failed:", slog.String("err", tx.Error.Error()))
		return tx.Error
	}

	return nil
}
