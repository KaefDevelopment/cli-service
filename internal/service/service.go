//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=interfaces_mock.go

package service

import (
	"context"

	"github.com/KaefDevelopment/cli-service/internal/model"
)

type CLIService struct {
	repo     Repository
	httpAddr string
	authKey  string
	txp      TxProvider
}

func NewCLIService(repo Repository, txp TxProvider, httpAddr, authKey string) *CLIService {
	return &CLIService{repo: repo, txp: txp, httpAddr: httpAddr, authKey: authKey}
}

type Repository interface {
	Create(ctx context.Context, events model.Events) error
	GetMarked(ctx context.Context) (model.Events, error)
	MarkSent(ctx context.Context) error
	Delete(ctx context.Context, events model.Events) error
	WithTx(txProvider TxProvider) Repository
}

type TxProvider interface {
	Transaction(f func(txProvider TxProvider) error) error
}
