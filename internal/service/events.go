package service

import (
	"context"
	"log"

	"github.com/KaefDevelopment/cli-service/internal/model"
	"github.com/KaefDevelopment/cli-service/internal/utils"
)

func (s *CLIService) CreateEvents(ctx context.Context, events model.Events) error {
	if err := s.repo.Create(ctx, events); err != nil {
		utils.WriteErrorResponse(err)
		return err
	}

	log.Printf("save %d events", len(events.Events))
	utils.WriteSuccessResponse()

	return nil
}
