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

	"gitlab.com/stash-password-manager/stash-server/config"
	"gitlab.com/stash-password-manager/stash-server/data"
	"gitlab.com/stash-password-manager/stash-server/migrations"
	"gitlab.com/stash-password-manager/stash-server/server"
)

func Run(configFilePath string) {
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

	// init sqlite connection
	err = data.InitSQLite(context.Background(), conf.SQLiteConfig.DbPath, migrations.Migrations())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initiate db connection: %s\n", err.Error())
		os.Exit(1)
	}

	// running http server
	serverOptions := server.ServerOptions{Addr: conf.Addr}
	httpServer := server.New(serverOptions)
	go func() {
		slog.Info("application server started", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	// graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh
	slog.Info("shutting down application")

	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(timeout); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
