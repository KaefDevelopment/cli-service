package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jaroslav1991/cli-service/internal/model"
	"github.com/jaroslav1991/cli-service/internal/service/dto"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime"
)

var (
	errBadStatusCode = errors.New("bad status code")
)

func (s *CLIService) Send(events model.Events, version, authKey string) error {
	if len(events.Events) == 0 {
		slog.Warn("empty events to send")
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
		}

		resEvent.Events = append(resEvent.Events, dtoEvent)
	}

	bytesEventsSend, err := json.Marshal(resEvent)
	if err != nil {
		slog.Error("fail marshal to sending:", slog.String("err", err.Error()))
		return err
	}

	req, err := http.NewRequest("POST", s.httpAddr, bytes.NewBuffer(bytesEventsSend))
	if err != nil {
		slog.Error("fail to send events:", slog.String("err", err.Error()))
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authKey)

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

	log.Printf("%s sent %d events", events.Events[0].AuthKey, len(events.Events))

	return nil
}
