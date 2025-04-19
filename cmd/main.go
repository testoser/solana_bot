package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pararti/solana-botyara/internal/config"
	"github.com/pararti/solana-botyara/internal/solana"
	"github.com/pararti/solana-botyara/internal/solana/copytrading"
	"github.com/pararti/solana-botyara/internal/utils"
)

func main() {
	logger := utils.NewLogger()

	cfg, err := config.Load("config.yaml")
	if err != nil {
		logger.Fatal("Failed to load configuration", err)
	}

	client, err := solana.NewClient(cfg.Solana.Endpoint, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Solana client", err)
	}

	wallet, err := solana.NewWalletFromPrivateKey(cfg.Wallet.PrivateKey, logger)
	if err != nil {
		logger.Fatal("Failed to initialize wallet", err)
	}

	strategy := copytrading.NewStrategy(cfg.Strategy, logger)

	monitor := copytrading.NewMonitor(client, cfg.Monitoring.Wallets, logger)

	executor := copytrading.NewExecutor(client, wallet, strategy, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		logger.Info("Shutting down gracefully...")
		cancel()
	}()

	logger.Info("Starting Solana copytrading bot...")

	txCh := make(chan *solana.Transaction)
	go monitor.Start(ctx, txCh)

	for {
		select {
		case tx := <-txCh:
			go executor.ProcessTransaction(ctx, tx)
		case <-ctx.Done():
			logger.Info("Bot stopped")
			return
		}
	}
}
