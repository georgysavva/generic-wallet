package main

import (
	"coins_wallet/config"
	"coins_wallet/postgres"
	"coins_wallet/wallet"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
)

func main() {

	var configPath string
	flag.StringVar(&configPath, "config", "config.json", "Path to the configuration file")
	flag.Parse()

	conf, err := config.Parse(configPath)
	if err != nil {
		panic(err)
	}

	paymentsRepository, err := postgres.NewPaymentsRepository(conf.Postgres)
	if err != nil {
		panic(err)
	}
	accountsRepository, err := postgres.NewAccountsRepository(conf.Postgres)
	if err != nil {
		panic(err)
	}

	walletService := wallet.NewService(paymentsRepository, accountsRepository)
	mux := http.NewServeMux()
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	httpLogger := log.With(logger, "component", "http")
	mux.Handle("/wallet/v1/", wallet.MakeHandler(walletService, httpLogger))

	httpAddr := fmt.Sprintf(":%d", conf.Port)
	server := &http.Server{Addr: httpAddr, Handler: mux}
	logger.Log("transport", "http", "address", httpAddr, "msg", "listening")
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	signalCode := waitingForShutdown()
	logger.Log("Received shutdown signal", "code", signalCode)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(conf.ShutDownTimeout))
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Log("msg", "Graceful shutdown failed", "err", err)
	}
}

func waitingForShutdown() os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	signalCode := <-ch
	return signalCode
}
