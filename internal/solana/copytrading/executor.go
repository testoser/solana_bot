package copytrading

import (
	"context"

	"github.com/pararti/solana-botyara/internal/solana"
	"go.uber.org/zap"
)

// Executor executes trades based on the strategy
type Executor struct {
	client   *solana.Client
	wallet   *solana.Wallet
	strategy *Strategy
	logger   *zap.Logger
}

// NewExecutor creates a new trade executor
func NewExecutor(client *solana.Client, wallet *solana.Wallet, strategy *Strategy, logger *zap.Logger) *Executor {
	return &Executor{
		client:   client,
		wallet:   wallet,
		strategy: strategy,
		logger:   logger,
	}
}

// ProcessTransaction processes a transaction and executes a copy trade if needed
func (e *Executor) ProcessTransaction(ctx context.Context, tx *solana.Transaction) {
	// Check if we should copy this trade
	if !e.strategy.ShouldCopyTrade(tx) {
		return
	}

	e.logger.Info("Copying trade",
		zap.String("signature", tx.Signature),
		zap.String("fromToken", tx.FromToken),
		zap.String("toToken", tx.ToToken),
		zap.Uint64("fromAmount", tx.FromAmount))

	// Adjust the trade amount based on our strategy
	amount := e.strategy.AdjustTradeAmount(tx.FromAmount)

	// Create the swap transaction
	swapTx, err := e.wallet.CreateSwapTransaction(
		ctx,
		e.client,
		tx.FromToken,
		tx.ToToken,
		amount,
	)
	if err != nil {
		e.logger.Error("Failed to create swap transaction", zap.Error(err))
		return
	}

	// Sign the transaction
	err = e.wallet.Sign(swapTx)
	if err != nil {
		e.logger.Error("Failed to sign transaction", zap.Error(err))
		return
	}

	// Send the transaction
	signature, err := e.client.SendTransaction(ctx, swapTx)
	if err != nil {
		e.logger.Error("Failed to send transaction", zap.Error(err))
		return
	}

	e.logger.Info("Trade executed successfully",
		zap.String("signature", signature),
		zap.String("fromToken", tx.FromToken),
		zap.String("toToken", tx.ToToken),
		zap.Uint64("amount", amount))
}
