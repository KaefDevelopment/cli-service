package repository

import (
	"context"
	"gorm.io/gorm/clause"
	"log/slog"

	"github.com/KaefDevelopment/cli-service/internal/model"
	"github.com/KaefDevelopment/cli-service/internal/service"
	"gorm.io/gorm"
)

const (
	defaultLimit = 10000
	batchSize    = 900
)

type CLIRepository struct {
	db *gorm.DB
}

func NewCLIRepository(db *gorm.DB) *CLIRepository {
	return &CLIRepository{db: db}
}

func (repo *CLIRepository) Create(ctx context.Context, events model.Events) error {
	err := repo.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&events.Events, batchSize).Error
	if err != nil {
		slog.Error("error with create gorm model:", slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (repo *CLIRepository) MarkSent(ctx context.Context) error {
	if err := repo.db.WithContext(ctx).Exec(
		"UPDATE events SET send=1 WHERE id IN (SELECT id FROM events WHERE send=0 ORDER BY createdAt limit ?)",
		defaultLimit).
		Error; err != nil {
		slog.Error("fail with update gorm column:", slog.String("err", err.Error()))
		return err
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
