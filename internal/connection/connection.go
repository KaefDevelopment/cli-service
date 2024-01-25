package connection

import (
	"fmt"
	"github.com/KaefDevelopment/cli-service/internal/utils"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"log/slog"
	"os"
)

func OpenDB(newConfigPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(newConfigPath + string(os.PathSeparator) + "cli.db"))
	if err != nil {
		slog.Error("open db failed:", slog.String("err", err.Error()))
		utils.WriteErrorResponse(utils.ErrConnectDB)
		return nil, fmt.Errorf("%w: %v", utils.ErrConnectDB, err)
	}

	return db, nil
}
