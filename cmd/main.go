package main

import (
	"cli-service/internal/connection"
	"cli-service/internal/model"
	cliservice "cli-service/internal/service"
	"cli-service/internal/service/repository"
	"flag"
	"fmt"
	"strings"
	"time"
)

var (
	inputData = flag.String(
		"d",
		"",
		"Stats data in JSON format string",
	)
)

func main() {
	flag.Parse()

	if strings.TrimSpace(*inputData) == "" {
		flag.Usage()
		return
	}

	now := time.Now()

	db, err := connection.OpenDB()
	if err != nil {
		return
	}

	if err := model.CreateTable(db); err != nil {
		return
	}

	repo := repository.NewCLIRepository(db)
	service := cliservice.NewCLIService(repo)
	data, err := service.ReadRequestData(*inputData)
	if err != nil {
		return
	}

	if err := service.Aggregate(data); err != nil {
		return
	}

	_, err = service.CreateEvents(data)
	if err != nil {
		return
	}

	getEvents, err := service.GetEvents()

	fmt.Println("get events:", getEvents)

	if err := service.Delete(); err != nil {
		return
	}

	fmt.Println("end time:", time.Since(now))

}
