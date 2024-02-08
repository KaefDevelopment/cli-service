package repository

import (
	"context"
	"log/slog"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/KaefDevelopment/cli-service/internal/model"
	"github.com/KaefDevelopment/cli-service/internal/service"
)

const defaultLimit = 10000

type CLIRepository struct {
	db *gorm.DB
}

func NewCLIRepository(db *gorm.DB) *CLIRepository {
	return &CLIRepository{db: db}
}

func (repo *CLIRepository) Create(ctx context.Context, events model.Events) error {
	err := repo.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&events.Events).Error
	if err != nil {
		slog.Error("error with create gorm model:", slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (repo *CLIRepository) MarkSent(ctx context.Context) error {
	err := repo.db.WithContext(ctx).Table("events").
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

func (repo *CLIRepository) GetMarked(ctx context.Context) (model.Events, error) {
	var events model.Events
	err := repo.db.WithContext(ctx).
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

func (repo *CLIRepository) Delete(ctx context.Context, events model.Events) error {
	err := repo.db.WithContext(ctx).Delete(events.Events).Error
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
