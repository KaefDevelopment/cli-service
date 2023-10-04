package service

import (
	"bytes"
	"cli-service/internal/model"
	"cli-service/internal/service/dto"
	"encoding/json"
	"log"
	"net/http"
)

func (s *CLIService) Send(events model.Events) error {
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
		log.Println("fail marshal to sending:", err)
		return err
	}

	for _, event := range events.Events {
		req, err := http.NewRequest("POST", s.httpAddr, bytes.NewBuffer(bytesEventsSend))
		if err != nil {
			log.Println("fail to send events:", err)
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", event.AuthKey)

		//fmt.Println("header:", req.Header.Get("Authorization"))

		client := http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			log.Println("fail with do sends:", err)
		}

		return resp.Body.Close()
	}

	return nil
}
