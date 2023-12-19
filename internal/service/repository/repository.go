package repository

import (
	"log/slog"

	"gorm.io/gorm"

	"github.com/jaroslav1991/cli-service/internal/model"
	"github.com/jaroslav1991/cli-service/internal/service"
)

const defaultLimit = 10000

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

func (repo *CLIRepository) MarkSent() error {
	err := repo.db.Table("events").
		Where("send = ?", 0).
		Order("createdAt asc").
		Limit(defaultLimit).
		Updates(map[string]interface{}{"send": 1}).
		Error
	if err != nil {
		slog.Error("fail with update gorm column:", slog.String("err", err.Error()))
	}

	return nil
}

func (repo *CLIRepository) GetMarked() (model.Events, error) {
	var events model.Events
	err := repo.db.
		Where("send = 1").
		Find(&events.Events).
		Limit(defaultLimit).
		Error
	if err != nil {
		slog.Error("fail with find events gorm:", slog.String("err", err.Error()))
		return model.Events{}, err
	}

	return events, nil
}

func (repo *CLIRepository) Delete(events model.Events) error {
	err := repo.db.Delete(events.Events).Error
	if err != nil {
		slog.Error("fail delete events gorm:", slog.String("err", err.Error()))
		return err
	}

	slog.Info("deleting successful")

	return nil
}

func (repo *CLIRepository) WithTx(txProvider service.TxProvider) service.Repository {
	return NewCLIRepository(txProvider.(*TxProvider).getTx())
}
