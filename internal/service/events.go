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

	// The stdout success/revoked marker is emitted once by the caller after the
	// whole event flow (persist + send) so plugins always receive a single JSON
	// object on stdout. Emitting it here too would produce a second line and
	// break the plugins that decode stdout as one value.
	return nil
}
