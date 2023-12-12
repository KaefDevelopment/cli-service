package repository

import (
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
	err := repo.db.Create(&events.Events).Error
	if err != nil {
		slog.Error("error with create gorm model:", slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (repo *CLIRepository) Update() error {
	err := repo.db.Table("events").Where("send = ?", 0).Limit(10000).Updates(map[string]interface{}{"send": 1}).Error
	if err != nil {
		slog.Error("fail with update gorm column:", slog.String("err", err.Error()))
	}

	return nil
}

func (repo *CLIRepository) Get() (model.Events, error) {
	var events model.Events

	err := repo.db.Find(&events.Events).Limit(10000).Error
	if err != nil {
		slog.Error("fail with find events gorm:", slog.String("err", err.Error()))
		return model.Events{}, err
	}

	return events, nil
}

func (repo *CLIRepository) Drop() error {
	err := repo.db.Delete(model.Event{}, "send =?", 1).Error
	if err != nil {
		slog.Error("fail delete events gorm:", slog.String("err", err.Error()))
		return err
	}

	slog.Info("deleting successful")

	return nil
}
