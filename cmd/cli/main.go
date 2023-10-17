package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jaroslav1991/cli-service/internal/connection"
	"github.com/jaroslav1991/cli-service/internal/model"
	cliservice "github.com/jaroslav1991/cli-service/internal/service"
	"github.com/jaroslav1991/cli-service/internal/service/repository"
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		return
	}

	fileInfo, err := os.OpenFile(homeDir+string(os.PathSeparator)+"cli-logger.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
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

	service.Aggregate(requestData)

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
