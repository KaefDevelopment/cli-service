package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/jaroslav1991/cli-service/internal/model"
	"github.com/jaroslav1991/cli-service/internal/service/dto"
)

var (
	errBadStatusCode = errors.New("bad status code")
)

func (s *CLIService) Send(version string) error {
	err := s.txp.Transaction(func(txp TxProvider) error {
		r := s.repo.WithTx(txp)

		events, err := s.lockEvents(r)
		if err != nil {
			return fmt.Errorf("failed to lock events: %w", err)
		}

		if err := s.sendEvents(events, version); err != nil {
			return fmt.Errorf("failed to sendEvents events: %w", err)
		}

		if err := s.unlockEvents(r, events); err != nil {
			return fmt.Errorf("failed to unlock events: %w", err)
		}

		return nil
	})

	if err != nil {
		slog.Error("transaction failed", slog.String("err", err.Error()))
		return err
	}
	return nil
}

func (s *CLIService) sendEvents(events model.Events, version string) error {
	if len(events.Events) == 0 {
		slog.Warn("empty events to sendEvents")
		return nil
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.httpAddr, bytes.NewBuffer(bytesEventsSend))
	if err != nil {
		slog.Error("fail to sendEvents events:", slog.String("err", err.Error()))
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
		slog.Error("fail status code", slog.String("status", resp.Status))
		return fmt.Errorf("%w: %s", errBadStatusCode, resp.Status)
	}

	log.Printf("%s sent %d events", s.authKey, len(events.Events))

	return nil
}

func (s *CLIService) lockEvents(repo Repository) (model.Events, error) {
	if err := repo.MarkSent(); err != nil {
		return model.Events{}, fmt.Errorf("failed to mark events: %w", err)
	}

	events, err := repo.GetMarked()
	if err != nil {
		return model.Events{}, fmt.Errorf("failed to get marked events: %w", err)
	}
	return events, nil
}

func (s *CLIService) unlockEvents(repo Repository, events model.Events) error {
	if err := repo.Delete(events); err != nil {
		return fmt.Errorf("failed to delete events: %w", err)
	}
	return nil
}
