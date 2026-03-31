package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/okunix/stash-server/cmd/stash-server/config"
	"github.com/okunix/stash-server/internal/adapter/cache"
	"github.com/okunix/stash-server/internal/adapter/postgres"
	"github.com/okunix/stash-server/internal/adapter/web"
	"github.com/okunix/stash-server/internal/core/services"
	"github.com/okunix/stash-server/migrations"
)

func Run(configFilePath string) {
	ctx := context.Background()
	// reading config file
	conf, err := config.ReadFromFile(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read config file: %s\n", err.Error())
		os.Exit(1)
	}

	// setting up slog
	logFilePath := conf.LogFile
	logFile := os.Stdout
	if logFilePath != "" {
		var err error
		logFile, err = os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to either create or open log file: %s\n", err.Error())
			os.Exit(1)
		}
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(logFile, nil)))

	// init postgres connection
	postgresInitParams := postgres.PostgresInitParams{
		User:       conf.PostgresConfig.User,
		Password:   conf.PostgresConfig.Password,
		Host:       conf.PostgresConfig.Host,
		Port:       conf.PostgresConfig.Port,
		SSLMode:    conf.PostgresConfig.SSLMode,
		Database:   conf.PostgresConfig.Database,
		Migrations: migrations.Migrations(),
	}
	err = postgres.Init(ctx, postgresInitParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initiate db connection: %s\n", err.Error())
		os.Exit(1)
	}

	// initializing user service instance
	userRepository := postgres.NewUserRepository(postgres.Postgres())
	userService := services.NewUserService(
		services.UserServiceParams{
			UserRepository: userRepository,
		},
	)

	// initializing admin user
	createdAdminUser, err := userService.InitializeAdminUser(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize admin user: %s\n", err.Error())
		os.Exit(1)
	}
	if createdAdminUser != nil {
		slog.Warn("created admin user",
			"username", createdAdminUser.Username,
			"password", createdAdminUser.Password,
		)
	}

	// initializing cache storage
	ephemeralStashStorage, err := cache.NewCache(5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize cache storage: %s\n", err.Error())
		os.Exit(1)
	}
	ephemeralUserIndex, err := cache.NewCache(5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize cache storage: %s\n", err.Error())
		os.Exit(1)
	}

	// initializing stash service instance
	stashRepository := postgres.NewStashRepository(postgres.Postgres())
	secretRepository := cache.NewSecretRepository(ephemeralStashStorage, ephemeralUserIndex)
	stashService := services.NewStashService(
		services.StashServiceParams{
			StashRepository:  stashRepository,
			SecretRepository: secretRepository,
		},
	)

	// running http server
	serverOptions := web.ServerOptions{
		Addr:         conf.Addr,
		UserService:  userService,
		StashService: stashService,
	}
	server := web.NewServer(serverOptions)
	go func() {
		slog.Info("application server started", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	// graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh
	slog.Info("shutting down application")

	timeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := server.Shutdown(timeout); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	postgres.Postgres().Close()
}
