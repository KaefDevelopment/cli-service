package service

import (
	"log"

	"github.com/KaefDevelopment/cli-service/internal/model"
	"github.com/KaefDevelopment/cli-service/internal/utils"
)

func (s *CLIService) CreateEvents(events model.Events) error {
	if err := s.repo.Create(events); err != nil {
		utils.WriteErrorResponse(err)
		return err
	}

	log.Printf("save %d events", len(events.Events))
	utils.WriteSuccessResponse()

	return nil
}
