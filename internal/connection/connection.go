package connection

import (
	"github.com/jaroslav1991/cli-service/internal/utils"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"log/slog"
	"os"
)

func OpenDB() (*gorm.DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("fail with home directory:", slog.String("err", err.Error()))
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(homeDir + string(os.PathSeparator) + "nau" + string(os.PathSeparator) + "cli.db"))
	if err != nil {
		slog.Error("open db failed:", slog.String("err", err.Error()))
		utils.WriteErrorResponse(utils.ErrConnectDB)
		return nil, err
	}

	return db, nil
}
