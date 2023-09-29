package service

import "cli-service/internal/model"

func (s *CLIService) CreateEvents(events model.Events) (model.Events, error) {
	result, err := s.repo.Create(events)
	if err != nil {
		return model.Events{}, err
	}

	return result, nil
}
