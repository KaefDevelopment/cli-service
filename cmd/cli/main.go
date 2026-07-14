package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"gorm.io/gorm/logger"

	"github.com/KaefDevelopment/cli-service/internal/connection"
	"github.com/KaefDevelopment/cli-service/internal/model"
	cliservice "github.com/KaefDevelopment/cli-service/internal/service"
	customlogger "github.com/KaefDevelopment/cli-service/internal/service/custom-logger"
	"github.com/KaefDevelopment/cli-service/internal/service/repository"
	"github.com/KaefDevelopment/cli-service/internal/utils"
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

	// wipeCmd lets a plugin explicitly clear the local offline queue, e.g. when it
	// learns about a REVOKED binding from its own status poll rather than from an
	// ingestion 410. It is idempotent: wiping an already-empty (or absent) queue
	// is harmless. Emits the same revoked marker on stdout on success.
	wipeCmd = &cobra.Command{
		Use:   "wipe",
		Short: "Wipe local plugin data (offline event queue)",
		Long:  "Drops all locally persisted events. Used by plugins to clear state on account deletion / unlink.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(cmd.Context(), time.Minute)
			defer cancel()

			homeDir, err := os.UserHomeDir()
			if err != nil {
				utils.WriteErrorResponse(err)
				log.Println(err)
				return
			}

			newConfigPath := homeDir + string(os.PathSeparator) + ".nau"
			oldConfigPath := homeDir + string(os.PathSeparator) + "nau"

			if err := utils.MigrateToNewConfigPath(newConfigPath, oldConfigPath); err != nil {
				utils.WriteErrorResponse(err)
				return
			}

			db, err := connection.OpenDB(logger.Default.LogMode(logger.Silent), newConfigPath)
			if err != nil {
				return
			}

			if err := model.InitSchema(ctx, db); err != nil {
				return
			}

			repo := repository.NewCLIRepository(db)
			if err := repo.WipeAll(ctx); err != nil {
				utils.WriteErrorResponse(err)
				return
			}

			utils.WriteRevokedResponse()
		},
	}

	eventCmd = &cobra.Command{
		Use:   "event",
		Short: "Event data in JSON format string",
		Long:  "using with event flag",
		Run: func(cmd *cobra.Command, args []string) {
			now := time.Now()
			ctx, cancel := context.WithTimeout(cmd.Context(), time.Minute)

			var err error
			defer func() {
				cancel()
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
				filepath.Join(newConfigPath, fmt.Sprintf("cli-logger-%s.txt", authKey[:8])),
				os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm,
			)
			if err != nil {
				log.Println(err)
				return
			}

			loggerSlog := slog.New(slog.NewTextHandler(fileInfo, &slog.HandlerOptions{Level: slog.LevelInfo}))
			slog.SetDefault(loggerSlog)

			newLogger := customlogger.New(loggerSlog)

			gormLogger := logger.New(
				newLogger,
				logger.Config{
					SlowThreshold:             0,
					Colorful:                  false,
					IgnoreRecordNotFoundError: false,
					ParameterizedQueries:      false,
					LogLevel:                  logger.Error,
				},
			)

			slog.Info("start cli...")

			db, err := connection.OpenDB(gormLogger, newConfigPath)
			if err != nil {
				return
			}

			if err := model.InitSchema(ctx, db); err != nil {
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

			if err = service.CreateEvents(ctx, requestData); err != nil {
				return
			}

			if authorized {
				if err := service.Send(ctx, version); err != nil {
					slog.Error("failed to send events", slog.String("err", err.Error()))

					// A REVOKED signal from the backend means this binding is gone.
					// Send() has already wiped the local queue; here we emit an
					// unambiguous machine-readable marker on stdout so the forking
					// plugin logs out and regenerates its uid. Clear the error so
					// deferred logging does not report it as a hard failure.
					if errors.Is(err, cliservice.ErrRevoked) {
						utils.WriteRevokedResponse()
						err = nil
					}
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

	rootCmd.AddCommand(eventCmd, versionCmd, wipeCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("can't execute root command:", err)
	}
}
