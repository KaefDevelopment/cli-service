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

	"github.com/go-git/go-git/v5"
	"github.com/google/uuid"

	"github.com/KaefDevelopment/cli-service/internal/model"
	"github.com/KaefDevelopment/cli-service/internal/service/dto"
)

const (
	repoInfo = "REPO_INFO"
	// revokedCode is the body marker the backend returns alongside HTTP 410 when
	// this plugin binding has been destroyed (account deleted or plugin unlinked).
	// See kaef-backend docs/plans/PLUGIN_DELETION_SIGNAL_PLAN.md §4.1.
	revokedCode = "PLUGIN_REVOKED"
)

var (
	errBadStatusCode = errors.New("bad status code")
	// ErrRevoked is a distinct, non-retryable error meaning the backend positively
	// declared this plugin binding gone (HTTP 410 + body code PLUGIN_REVOKED).
	// It is deliberately NOT folded into errBadStatusCode so the caller can wipe
	// local data instead of retrying the queue forever.
	ErrRevoked = errors.New("plugin revoked")
)

// revokedBody is the minimal shape of the backend ErrorDto we need to confirm a
// revoke; extra fields (e.g. description with the reason) are ignored on purpose.
type revokedBody struct {
	Code string `json:"code"`
}

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
		if errors.Is(err, ErrRevoked) {
			// The send transaction above has already rolled back, so the events
			// are still on disk. Wipe them via the base repo (s.repo, NOT the
			// rolled-back tx repo) so the delete commits on its own connection
			// and is not undone together with the failed send. We do NOT retry.
			slog.Warn("plugin revoked, wiping local data", slog.String("pluginId", s.authKey))

			if wipeErr := s.repo.WipeAll(ctx, s.authKey); wipeErr != nil {
				slog.Error("failed to wipe local data after revoke", slog.String("err", wipeErr.Error()))
			}

			return err
		}

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
		Events:     make([]dto.Event, 0, len(events.Events)+1),
	}

	var (
		repoURLs = make(map[string]struct{})
		info     = model.Params{
			"reposInfo": []map[string]interface{}{},
		}
	)

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

		if _, ok := repoURLs[events.Events[i].ProjectBaseDir]; !ok {
			repoURLs[events.Events[i].ProjectBaseDir] = struct{}{}

			m := map[string]interface{}{
				"repoUrls":       getURLsByDir(events.Events[i].ProjectBaseDir),
				"projectBaseDir": events.Events[i].ProjectBaseDir,
				"project":        events.Events[i].Project,
			}
			info["reposInfo"] = append(info["reposInfo"].([]map[string]interface{}), m)
		}

		resEvent.Events = append(resEvent.Events, dtoEvent)

	}

	repoUrl := dto.Event{
		Id:        uuid.New().String(),
		CreatedAt: time.Now().String(),
		Type:      repoInfo,
		Timezone:  events.Events[0].Timezone,
		PluginId:  events.Events[0].AuthKey,
		Params:    info,
	}
	resEvent.Events = append(resEvent.Events, repoUrl)

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

		// Distinguish a positive REVOKED signal from every other >=400. Only a
		// 410 whose body parses as {"code":"PLUGIN_REVOKED"} qualifies (both
		// checks required); anything else stays a generic keep-and-retry error.
		if resp.StatusCode == http.StatusGone {
			var parsed revokedBody
			if err := json.Unmarshal(body, &parsed); err == nil && parsed.Code == revokedCode {
				return ErrRevoked
			}
		}

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

func getURLsByDir(projectBaseDir string) []string {
	repo, err := git.PlainOpen(projectBaseDir)
	if err != nil {
		slog.Warn("fail to open repository:",
			slog.String("err", err.Error()),
			slog.String("projectBaseDir", projectBaseDir))
		return nil
	}

	remotes, err := repo.Remotes()
	if err != nil {
		slog.Warn("fail to get remotes:",
			slog.String("err", err.Error()),
			slog.String("projectBaseDir", projectBaseDir))
		return nil
	}

	var urls []string
	for _, remote := range remotes {
		urls = append(urls, remote.Config().URLs...)
	}

	return urls
}
