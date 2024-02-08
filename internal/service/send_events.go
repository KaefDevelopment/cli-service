package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/KaefDevelopment/cli-service/internal/model"
	"github.com/KaefDevelopment/cli-service/internal/service/dto"
)

var (
	errBadStatusCode = errors.New("bad status code")
)

func (s *CLIService) sendWithLock(ctx context.Context, r Repository, version string) error {
	events, err := s.lockEvents(ctx, r)
	if err != nil {
		return fmt.Errorf("failed to lock events: %w", err)
	}

	if len(events.Events) == 0 {
		slog.Warn("no events to send")
		return nil
	}

	if err := s.sendEvents(ctx, events, version); err != nil {
		return fmt.Errorf("failed to sendEvents events: %w", err)
	}

	if err := s.unlockEvents(ctx, r, events); err != nil {
		return fmt.Errorf("failed to unlock events: %w", err)
	}

	return nil
}

func (s *CLIService) Send(ctx context.Context, version string) error {
	err := s.txp.Transaction(func(txp TxProvider) error {
		return s.sendWithLock(ctx, s.repo.WithTx(txp), version)
	})

	if err != nil {
		slog.Error("transaction failed", slog.String("err", err.Error()))
		return err
	}
	return nil
}

func (s *CLIService) sendEvents(ctx context.Context, events model.Events, version string) error {
	osName := runtime.GOOS

	deviceName, _ := os.Hostname()

	resEvent := dto.SendEvents{
		OsName:     osName,
		DeviceName: deviceName,
		CliVersion: version,
		Events:     make([]dto.Event, 0, len(events.Events)),
	}

	for i := range events.Events {
		dtoEvent := dto.Event{
			Id:             events.Events[i].Id,
			CreatedAt:      events.Events[i].CreatedAt,
			Type:           events.Events[i].Type,
			Project:        events.Events[i].Project,
			ProjectBaseDir: events.Events[i].ProjectBaseDir,
			Language:       events.Events[i].Language,
			Target:         events.Events[i].Target,
			Branch:         events.Events[i].Branch,
			Timezone:       events.Events[i].Timezone,
			Params:         events.Events[i].Params,
			PluginId:       events.Events[i].AuthKey,
		}

		resEvent.Events = append(resEvent.Events, dtoEvent)
	}

	bytesEventsSend, err := json.Marshal(resEvent)
	if err != nil {
		slog.Error("fail marshal to sending:", slog.String("err", err.Error()))
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.httpAddr, bytes.NewBuffer(bytesEventsSend))
	if err != nil {
		slog.Error("failed to send events:", slog.String("err", err.Error()))
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.authKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Warn("fail with do sends:", slog.String("err", err.Error()))
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		slog.Error("fail status code", slog.String("status", resp.Status), slog.String("request", string(bytesEventsSend)), slog.String("response", string(body)))
		return fmt.Errorf("%w: %s", errBadStatusCode, resp.Status)
	}

	log.Printf("%s sent %d events", s.authKey, len(events.Events))

	return nil
}

func (s *CLIService) lockEvents(ctx context.Context, repo Repository) (model.Events, error) {
	if err := repo.MarkSent(ctx); err != nil {
		return model.Events{}, fmt.Errorf("failed to mark events: %w", err)
	}

	events, err := repo.GetMarked(ctx)
	if err != nil {
		return model.Events{}, fmt.Errorf("failed to get marked events: %w", err)
	}
	return events, nil
}

func (s *CLIService) unlockEvents(ctx context.Context, repo Repository, events model.Events) error {
	if err := repo.Delete(ctx, events); err != nil {
		return fmt.Errorf("failed to delete events: %w", err)
	}
	return nil
}
