package connection

import (
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"log"
)

func OpenDB() (*gorm.DB, error) {

	db, err := gorm.Open(sqlite.Open("./cli.db"))
	if err != nil {
		log.Println("open db failed:", err)
		return nil, err
	}

	return db, nil
}
