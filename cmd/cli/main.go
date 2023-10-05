package main

import (
	"cli-service/internal/connection"
	"cli-service/internal/model"
	cliservice "cli-service/internal/service"
	"cli-service/internal/service/repository"
	"flag"
	"log"
	"os"
	"strings"
	"time"
)

var (
	inputData = flag.String(
		"d",
		"",
		"Request data in JSON format string",
	)

	httpAddr = flag.String(
		"s",
		"http://localhost:8181/events",
		"Http address for sending events",
	)

	authKey = flag.String(
		"k",
		"",
		"authorization key",
	)
)

func init() {
	fileInfo, err := os.OpenFile("cli-logger.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}
	log.SetOutput(fileInfo)

}

func main() {
	log.Println("start cli...")

	flag.Parse()

	now := time.Now()

	var err error
	defer func() {
		if err != nil {
			log.Println("error:", err)
		}

		log.Println("ending time:", time.Since(now))
		log.Println("ending cli")
	}()

	if strings.TrimSpace(*inputData) == "" {
		flag.Usage()
		return
	}

	db, err := connection.OpenDB()
	if err != nil {
		return
	}

	if err := model.CreateTable(db); err != nil {
		return
	}

	repo := repository.NewCLIRepository(db)
	service := cliservice.NewCLIService(repo, *httpAddr, *authKey)

	requestData, err := service.ReadRequestData(*inputData)
	if err != nil {
		return
	}

	if err := service.Aggregate(requestData); err != nil {
		return
	}

	if err = service.CreateEvents(requestData); err != nil {
		return
	}

	if err = service.UpdateEvents(); err != nil {
		return
	}

	keys, err := service.GetKeys()
	if err != nil {
		return
	}

	eventsToSend, err := service.GetEvents(keys)

	for _, event := range eventsToSend.Events {
		if err := service.Send(event); err != nil {
			return
		}
	}

	if err := service.Delete(); err != nil {
		return
	}
}
