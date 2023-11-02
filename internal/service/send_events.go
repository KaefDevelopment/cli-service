package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"

	"github.com/jaroslav1991/cli-service/internal/model"
	"github.com/jaroslav1991/cli-service/internal/service/dto"
)

var (
	errInternalServer = errors.New("internal server error")
)

func (s *CLIService) Send(events model.Events) error {
	if len(events.Events) == 0 {
		slog.Warn("empty events to send")
		return nil
	}

	var resEvent dto.Events

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
		}

		resEvent.Events = append(resEvent.Events, dtoEvent)
	}

	bytesEventsSend, err := json.Marshal(resEvent)
	if err != nil {
		slog.Error("fail marshal to sending:", err)
		return err
	}

	req, err := http.NewRequest("POST", s.httpAddr, bytes.NewBuffer(bytesEventsSend))
	if err != nil {
		slog.Error("fail to send events:", slog.String("err", err.Error()))
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", events.Events[0].AuthKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Warn("fail with do sends:", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		slog.Error("fail status code:", slog.Any("err", resp.Header))
		return errInternalServer
	}

	log.Printf("sent %d events", len(events.Events))

	return nil
}
