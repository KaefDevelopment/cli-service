package service

import (
	"encoding/json"
	"fmt"
	"github.com/jaroslav1991/cli-service/internal/model"
	"github.com/jaroslav1991/cli-service/internal/utils"
	"log"
)

func (s *CLIService) ReadRequestData(request string) (model.Events, error) {
	var requestModel model.Events

	if err := json.Unmarshal([]byte(request), &requestModel); err != nil {
		log.Println("read data unmarshal failed:", err)
		utils.WriteErrorResponse(utils.ErrReadRequestDataUnmarshal)
		return model.Events{}, err
	}

	for i := range requestModel.Events {
		requestModel.Events[i].AuthKey = s.authKey

		if s.authKey == "" {
			log.Println("fail with auth key, couldn't be empty")
			utils.WriteErrorResponse(utils.ErrAuthKey)
			return model.Events{}, fmt.Errorf("%v", utils.ErrAuthKey)
		}
	}

	utils.WriteSuccessResponse(utils.ResponseForPlugin{
		Status: "Success",
	})

	return requestModel, nil
}
