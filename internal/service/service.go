//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=interfaces_mock.go

package service

import (
	"github.com/jaroslav1991/cli-service/internal/model"
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
	Create(events model.Events) error
	GetMarked() (model.Events, error)
	MarkSent() error
	Delete(events model.Events) error
	WithTx(txProvider TxProvider) Repository
}

type TxProvider interface {
	Transaction(f func(txProvider TxProvider) error) error
}
