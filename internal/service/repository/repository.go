package repository

import (
	"context"
	"log/slog"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/KaefDevelopment/cli-service/internal/model"
	"github.com/KaefDevelopment/cli-service/internal/service"
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

// WipeAll removes the locally persisted plugin data for authKey. Today the queue
// lives in a single table (events); when agent tables land (agent_cursors/
// agent_state) they should be wiped here too, scoped by the same key. The delete
// is scoped by authKey so a revoke of one plugin does not drop another plugin's
// queued events from the shared ~/.nau/cli.db. It runs on its own connection/
// transaction so a caller can invoke it after a failed, rolled-back send without
// the delete being rolled back along with the send.
func (repo *CLIRepository) WipeAll(ctx context.Context, authKey string) error {
	if err := repo.db.WithContext(ctx).Exec("DELETE FROM events WHERE authKey = ?", authKey).Error; err != nil {
		slog.Error("fail wipe events gorm:", slog.String("err", err.Error()))
		return err
	}

	slog.Info("wiped local events after revoke")

	return nil
}

func (repo *CLIRepository) WithTx(txProvider service.TxProvider) service.Repository {
	return NewCLIRepository(txProvider.(*TxProvider).getTx())
}
