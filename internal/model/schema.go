package model

import (
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"github.com/jaroslav1991/cli-service/internal/utils"
)

const (
	initSchema = `create table if not exists events (
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
	);
	create unique index if not exists pk_id on events (id);
	create index if not exists send_created_at_idx on events (send, createdAt);`
)

func InitSchema(db *gorm.DB) error {
	err := db.Exec(initSchema).Error
	if err != nil {
		slog.Error("create table failed:", slog.String("err", err.Error()))
		utils.WriteErrorResponse(utils.ErrCreateTable)
		return fmt.Errorf("%w: %v", utils.ErrCreateTable, err)
	}

	return nil
}
