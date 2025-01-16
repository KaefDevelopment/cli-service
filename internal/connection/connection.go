package connection

import (
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/KaefDevelopment/cli-service/internal/utils"
)

func OpenDB(logger logger.Interface, newConfigPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(newConfigPath+string(os.PathSeparator)+"cli.db"), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		slog.Error("open db failed:", slog.String("err", err.Error()))
		utils.WriteErrorResponse(utils.ErrConnectDB)
		return nil, fmt.Errorf("%w: %v", utils.ErrConnectDB, err)
	}

	//_ = db.Exec("PRAGMA journal_mode = WAL") // set journal mode = WAL (Write Ahead Logging)

	return db, nil
}
