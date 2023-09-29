package service

import (
	"cli-service/internal/model"
)

type CLIService struct {
	repo Repository
}

func NewCLIService(repo Repository) *CLIService {
	return &CLIService{repo: repo}
}

type Repository interface {
	Create(events model.Events) (model.Events, error)
	Get() (model.Events, error)
	Drop() error
}
