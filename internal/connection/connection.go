package connection

import (
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"log"
	"os"
)

func OpenDB() (*gorm.DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("fail with home directory:", err)
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(homeDir + string(os.PathSeparator) + "cli.db"))
	if err != nil {
		log.Println("open db failed:", err)
		return nil, err
	}

	return db, nil
}
