package service

import (
	"cli-service/internal/model"
	"encoding/json"
	"fmt"
	"log"
)

func (s *CLIService) ReadRequestData(request string) (model.Events, error) {
	var requestModel model.Events

	if err := json.Unmarshal([]byte(request), &requestModel); err != nil {
		log.Println("read data unmarshal failed:", err)
		return model.Events{}, err
	}

	fmt.Println("len events:", len(requestModel.Events))
	return requestModel, nil
}
