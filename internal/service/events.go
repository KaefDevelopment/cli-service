package service

import (
	"github.com/jaroslav1991/cli-service/internal/model"
)

func (s *CLIService) CreateEvents(events model.Events) error {
	return s.repo.Create(events)
}

func (s *CLIService) UpdateEvents() error {
	return s.repo.Update()
}

func (s *CLIService) GetKeys() ([]string, error) {
	keys, err := s.repo.GetAuthKeys()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (s *CLIService) GetEvents(keys []string) (model.EventsByAuthKey, error) {
	events, err := s.repo.Get(keys)
	if err != nil {
		return model.EventsByAuthKey{}, err
	}

	return events, nil
}

func (s *CLIService) Delete() error {
	return s.repo.Drop()
}
