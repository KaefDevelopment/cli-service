package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/jaroslav1991/cli-service/internal/model"
)

var (
	errAuthKey = errors.New("failed with auth key from request events")
)

func (s *CLIService) ReadRequestData(request string) (model.Events, error) {
	var requestModel model.Events

	if err := json.Unmarshal([]byte(request), &requestModel); err != nil {
		log.Println("read data unmarshal failed:", err)
		return model.Events{}, err
	}

	for i := range requestModel.Events {
		requestModel.Events[i].AuthKey = s.authKey

		if s.authKey == "" {
			log.Println("fail with auth key, couldn't be empty")
			return model.Events{}, fmt.Errorf("%v", errAuthKey)
		}
	}

	return requestModel, nil
}
