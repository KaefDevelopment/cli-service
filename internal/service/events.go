package service

import (
	"github.com/jaroslav1991/cli-service/internal/model"
	"github.com/jaroslav1991/cli-service/internal/utils"
	"log"
)

func (s *CLIService) CreateEvents(events model.Events) error {
	if err := s.repo.Create(events); err != nil {
		utils.WriteErrorResponse(err)
		return err
	}

	log.Printf("save %d events", len(events.Events))
	utils.WriteSuccessResponse(utils.ResponseForPlugin{
		Status: true,
	})

	return nil
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
