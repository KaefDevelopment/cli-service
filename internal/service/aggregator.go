package service

import (
	"log"
	"os"
	"strings"

	"github.com/google/uuid"

	"github.com/jaroslav1991/cli-service/internal/model"
)

var (
	getBranchFn = GetBranchByProjectBaseDir
	getIDFn     = GetUUID
)

func GetBranchByProjectBaseDir(projectBaseDir string) string {
	filename := projectBaseDir + string(os.PathSeparator) + ".git" + string(os.PathSeparator) + "HEAD"

	currentBranch, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("current branch path not found: %v", err)
		return ""
	}

	return strings.TrimSpace(strings.ReplaceAll(string(currentBranch), "ref: refs/heads/", ""))
}

func GetUUID() string {
	return uuid.New().String()
}

func (s *CLIService) Aggregate(events model.Events) error {
	for i := range events.Events {
		if events.Events[i].Branch != "" {
			continue
		}

		if eventBranch := getBranchFn(events.Events[i].ProjectBaseDir); eventBranch != "" {
			events.Events[i].Branch = eventBranch
		}

		if events.Events[i].Id != "" {
			continue
		} else if events.Events[i].Id == "" {
			events.Events[i].Id = getIDFn()
		}

	}

	return nil
}
