package service

import (
	"cli-service/internal/model"
)

func (s *CLIService) CreateEvents(events model.Events) error {
	if err := s.repo.Create(events); err != nil {
		return err
	}

	return nil
}

func (s *CLIService) UpdateEvents() error {
	if err := s.repo.Update(); err != nil {
		return err
	}

	return nil
}

func (s *CLIService) GetKeys() ([]string, error) {
	keys, err := s.repo.GetAuthKeys()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (s *CLIService) GetEvents(keys []string) (model.Response, error) {
	events, err := s.repo.Get(keys)
	if err != nil {
		return model.Response{}, err
	}

	return events, nil
}

func (s *CLIService) Delete() error {
	if err := s.repo.Drop(); err != nil {
		return err
	}

	return nil
}
