package service

import (
	"log/slog"
	"os"
	"strings"

	"github.com/google/uuid"

	"github.com/jaroslav1991/cli-service/internal/model"
)

var (
	getBranchFn = GetBranchByProjectBaseDir
	getIDFn     = GetUUID
	branchCache = make(map[string]string)
)

func GetBranchByProjectBaseDir(projectBaseDir string) string {
	if cachedBranch, ok := branchCache[projectBaseDir]; ok {
		return cachedBranch
	}

	filename := projectBaseDir + string(os.PathSeparator) + ".git" + string(os.PathSeparator) + "HEAD"

	currentBranch, err := os.ReadFile(filename)
	if err != nil {
		slog.Warn("current branch path not found:", slog.String("err", err.Error()))
		branchCache[projectBaseDir] = ""
	}

	branchName := strings.TrimSpace(strings.ReplaceAll(string(currentBranch), "ref: refs/heads/", ""))
	branchCache[projectBaseDir] = branchName

	return branchName
}

func GetUUID() string {
	return uuid.New().String()
}

func (s *CLIService) Aggregate(events model.Events) {
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
}
