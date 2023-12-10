package service

import (
	"encoding/json"
	"fmt"
	"github.com/jaroslav1991/cli-service/internal/model"
	"github.com/jaroslav1991/cli-service/internal/utils"
	"log/slog"
)

func (s *CLIService) ReadRequestData(request string) (model.Events, error) {
	var requestModel model.Events

	if err := json.Unmarshal([]byte(request), &requestModel); err != nil {
		slog.Error("read data unmarshal failed:", slog.String("json", request), slog.String("err", err.Error()))
		utils.WriteErrorResponse(utils.ErrReadRequestDataUnmarshal)
		return model.Events{}, fmt.Errorf("%w:%v", utils.ErrReadRequestDataUnmarshal, err)
	}

	for i := range requestModel.Events {
		requestModel.Events[i].AuthKey = s.authKey

		if s.authKey == "" {
			slog.Error("fail with auth key, couldn't be empty", slog.String("err", utils.ErrAuthKey.Error()))
			utils.WriteErrorResponse(utils.ErrAuthKey)
			return model.Events{}, fmt.Errorf("%v", utils.ErrAuthKey)
		}
	}

	return requestModel, nil
}
