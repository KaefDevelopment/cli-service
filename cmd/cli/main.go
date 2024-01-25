package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/KaefDevelopment/cli-service/internal/connection"
	"github.com/KaefDevelopment/cli-service/internal/model"
	cliservice "github.com/KaefDevelopment/cli-service/internal/service"
	"github.com/KaefDevelopment/cli-service/internal/service/repository"
	"github.com/KaefDevelopment/cli-service/internal/utils"

	"github.com/spf13/cobra"
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
		Short: "GetMarked cli version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	eventCmd = &cobra.Command{
		Use:   "event",
		Short: "Event data in JSON format string",
		Long:  "using with event flag",
		Run: func(cmd *cobra.Command, args []string) {

			now := time.Now()

			var err error
			defer func() {
				if err != nil {
					slog.Error("error:", slog.String("err", err.Error()))
				}

				log.Println("ending time:", time.Since(now))
				log.Println("ending cli")
			}()

			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Println(err)
				return
			}

			newConfigPath := homeDir + string(os.PathSeparator) + ".nau"
			oldConfigPath := homeDir + string(os.PathSeparator) + "nau"

			if err := utils.MigrateToNewConfigPath(newConfigPath, oldConfigPath); err != nil {
				return
			}

			fileInfo, err := os.OpenFile(
				filepath.Join(newConfigPath, fmt.Sprintf("cli-logger-%s.txt", authKey)),
				os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm,
			)
			if err != nil {
				log.Println(err)
				return
			}

			logger := slog.New(slog.NewTextHandler(fileInfo, &slog.HandlerOptions{Level: slog.LevelInfo}))
			slog.SetDefault(logger)

			slog.Info("start cli...")

			db, err := connection.OpenDB(newConfigPath)
			if err != nil {
				return
			}

			if err := model.InitSchema(db); err != nil {
				return
			}

			repo := repository.NewCLIRepository(db)
			txp := repository.NewTxProvider(db)
			service := cliservice.NewCLIService(repo, txp, httpAddr, authKey)

			requestData, err := service.ReadRequestData(inputData)
			if err != nil {
				return
			}

			service.Aggregate(requestData)

			if err = service.CreateEvents(requestData); err != nil {
				return
			}

			if authorized {
				if err := service.Send(version); err != nil {
					slog.Error("failed to send events", slog.String("err", err.Error()))
				}
			} else {
				slog.Warn("plugin is not authorized", slog.String("pluginId", authKey))
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
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("can't execute root command:", err)
	}
}
