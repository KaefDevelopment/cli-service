//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=interfaces_mock.go

package service

import (
	"github.com/jaroslav1991/cli-service/internal/model"
)

type CLIService struct {
	repo       Repository
	httpAddr   string
	authKey    string
	authorized bool
}

func NewCLIService(repo Repository, httpAddr, authKey string, authorized bool) *CLIService {
	return &CLIService{repo: repo, httpAddr: httpAddr, authKey: authKey, authorized: authorized}
}

type Repository interface {
	Create(events model.Events) error
	Get() (model.Events, error)
	Update() error
	Drop() error
}
