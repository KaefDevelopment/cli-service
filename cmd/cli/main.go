package main

import (
	"errors"
	"fmt"
	"github.com/jaroslav1991/cli-service/internal/connection"
	"github.com/jaroslav1991/cli-service/internal/model"
	cliservice "github.com/jaroslav1991/cli-service/internal/service"
	"github.com/jaroslav1991/cli-service/internal/service/repository"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"os"
	"time"
)

var (
	inputData  string
	httpAddr   string
	authKey    string
	authorized bool
	version    string

	rootCmd = &cobra.Command{
		Use:   "cli",
		Short: "Root command",
		Long:  "Root command for CLI",
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Get cli version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	eventCmd = &cobra.Command{
		Use:   "event",
		Short: "Event data in JSON format string",
		Long:  "using with event flag",
		Run: func(cmd *cobra.Command, args []string) {
			slog.Info("start cli...")

			now := time.Now()

			var err error
			defer func() {
				if err != nil {
					slog.Error("error:", slog.String("err", err.Error()))
				}

				log.Println("ending time:", time.Since(now))
				log.Println("ending cli")
			}()

			db, err := connection.OpenDB()
			if err != nil {
				return
			}

			if err := model.CreateTable(db); err != nil {
				return
			}

			repo := repository.NewCLIRepository(db)
			service := cliservice.NewCLIService(repo, httpAddr, authKey, authorized)

			requestData, err := service.ReadRequestData(inputData)
			if err != nil {
				return
			}

			service.Aggregate(requestData)

			if err = service.CreateEvents(requestData); err != nil {
				return
			}

			if authorized {

				keys, err := service.GetKeys()
				if err != nil {
					return
				}

				eventsToSend, err := service.GetEvents(keys)

				for _, event := range eventsToSend.Events {
					if err := service.Send(event, version, authKey); err != nil {
						return
					}
				}

				if err = service.UpdateEvents(); err != nil {
					return
				}

				if err := service.Delete(); err != nil {
					return
				}
			} else {
				slog.Warn("user not authorized")
				return
			}
		},
	}
)

func init() {
	eventCmd.Flags().StringVarP(&inputData, "data", "d", "", "Request data in JSON format string")
	eventCmd.Flags().StringVarP(&authKey, "auth-key", "k", "", "Authorization key")
	eventCmd.Flags().StringVarP(&httpAddr, "server-host", "s", "https://nautime.io/api/plugin/v1/events?source=cli&version=$version", "Http address for sending events")
	eventCmd.Flags().BoolVarP(&authorized, "authorized", "a", true, "Take info about authorization")

	rootCmd.AddCommand(eventCmd, versionCmd)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		return
	}

	if err := os.Mkdir(homeDir+string(os.PathSeparator)+"nau", os.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
		log.Println(err)
		return
	}

	fileInfo, err := os.OpenFile(
		homeDir+string(os.PathSeparator)+"nau"+string(os.PathSeparator)+"cli-logger.txt",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm,
	)
	if err != nil {
		log.Println(err)
		return
	}

	logger := slog.New(slog.NewTextHandler(fileInfo, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("can't execute root command:", err)
	}
}
