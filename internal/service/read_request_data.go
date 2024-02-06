package service

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/KaefDevelopment/cli-service/internal/model"
	"github.com/KaefDevelopment/cli-service/internal/utils"
)

func (s *CLIService) ReadRequestData(request string) (model.Events, error) {
	var requestModel model.Events

	if err := json.Unmarshal([]byte(request), &requestModel); err != nil {
		slog.Error("read data unmarshal failed:", slog.String("json", request), slog.String("err", err.Error()))
		utils.WriteErrorResponse(utils.ErrReadRequestDataUnmarshal)
		return model.Events{}, fmt.Errorf("%w:%v", utils.ErrReadRequestDataUnmarshal, err)
	}

	if s.authKey == "" {
		slog.Error("fail with auth key, couldn't be empty", slog.String("err", utils.ErrAuthKey.Error()))
		utils.WriteErrorResponse(utils.ErrAuthKey)
		return model.Events{}, fmt.Errorf("%v", utils.ErrAuthKey)
	}

	for i := range requestModel.Events {
		requestModel.Events[i].AuthKey = s.authKey

		if requestModel.Events[i].Type == "" {
			slog.Error("fail with type field, couldn't be empty", slog.String("err", utils.ErrTypeField.Error()))
			utils.WriteErrorResponse(utils.ErrTypeField)
			return model.Events{}, fmt.Errorf("%v", utils.ErrTypeField)
		}

		if requestModel.Events[i].CreatedAt == "" {
			slog.Error("fail with created at field, couldn't be empty", slog.String("err", utils.ErrCreatedAtField.Error()))
			utils.WriteErrorResponse(utils.ErrCreatedAtField)
			return model.Events{}, fmt.Errorf("%v", utils.ErrCreatedAtField)
		}
	}

	return requestModel, nil
}
