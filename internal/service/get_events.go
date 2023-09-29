package service

import "cli-service/internal/model"

func (s *CLIService) GetEvents() (model.Events, error) {
	events, err := s.repo.Get()
	if err != nil {
		return model.Events{}, err
	}

	return events, nil
}
