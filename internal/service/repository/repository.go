package repository

import (
	"log"
	"log/slog"

	"gorm.io/gorm"

	"github.com/jaroslav1991/cli-service/internal/model"
)

type CLIRepository struct {
	db *gorm.DB
}

func NewCLIRepository(db *gorm.DB) *CLIRepository {
	return &CLIRepository{db: db}
}

func (repo *CLIRepository) Create(events model.Events) error {
	result := repo.db.Create(&events.Events)
	if result.Error != nil {
		slog.Error("error with create gorm model:", slog.String("err", result.Error.Error()))
		return result.Error
	}

	return nil
}

func (repo *CLIRepository) Update() error {
	result := repo.db.Table("events").Where("send = ?", 0).Limit(10000).Updates(map[string]interface{}{"send": 1})
	if result.Error != nil {
		slog.Error("fail with update gorm column:", slog.String("err", result.Error.Error()))
	}

	return nil
}

func (repo *CLIRepository) Get(authKey []string) (model.EventsByAuthKey, error) {
	var res model.EventsByAuthKey

	for _, key := range authKey {
		var events model.Events

		result := repo.db.Where("send = ?", 1).Where("authKey = ?", key).Limit(10000).Find(&events.Events)
		if result.Error != nil {
			slog.Error("fail with find events gorm:", slog.String("err", result.Error.Error()))
			return model.EventsByAuthKey{}, result.Error
		}

		res.Events = append(res.Events, events)
	}

	return res, nil
}

func (repo *CLIRepository) GetAuthKeys() ([]string, error) {
	result := repo.db.Table("events")

	result = result.Raw("select distinct authKey from events")

	rows, err := result.Rows()
	if err != nil {
		slog.Error("fail with rows gorm:", slog.String("err", err.Error()))
		return nil, err
	}

	defer rows.Close()

	var keys []string

	for rows.Next() {
		var authKey model.Event

		if err := rows.Scan(&authKey.AuthKey); err != nil {
			slog.Error("fail rows scan gorm:", slog.String("err", err.Error()))
			return nil, err
		}

		keys = append(keys, authKey.AuthKey)
	}

	return keys, nil
}

func (repo *CLIRepository) Drop(events model.EventsByAuthKey) error {
	delCounter := 0

	for i := range events.Events {

		result := repo.db.Delete(&events.Events[i].Events, "send = ?", "1")
		if result.Error != nil {
			slog.Error("fail delete events gorm:", slog.String("err", result.Error.Error()))
			return result.Error
		}

		delCounter += len(events.Events[i].Events)
	}

	log.Printf("deleted %d events", delCounter)

	slog.Info("deleting successful")

	return nil
}
