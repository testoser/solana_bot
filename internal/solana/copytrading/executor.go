package copytrading

import (
	"context"
	"log"

	"github.com/pararti/solana-botyara/internal/solana"
)

// Executor executes trades based on the strategy
type Executor struct {
	client   *solana.Client
	wallet   *solana.Wallet
	strategy *Strategy
	logger   *log.Logger
}

// NewExecutor creates a new trade executor
func NewExecutor(client *solana.Client, wallet *solana.Wallet, strategy *Strategy, logger *log.Logger) *Executor {
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

	e.logger.Println("Copying trade",
		"signature", tx.Signature,
		"fromToken", tx.FromToken,
		"toToken", tx.ToToken,
		"fromAmount", tx.FromAmount)

	// Adjust the trade amount based on our strategy
	amount := e.strategy.AdjustTradeAmount(tx.FromAmount)

	// Create the swap transaction
	swapTx, err := e.wallet.CreateSwapTransaction(
		ctx,
		e.client,
	)
	if err != nil {
		e.logger.Println("Failed to create swap transaction", err)
		return
	}

	// Sign the transaction
	err = e.wallet.Sign(swapTx)
	if err != nil {
		e.logger.Println("Failed to sign transaction")
		return
	}

	// Send the transaction
	signature, err := e.client.SendTransaction(ctx, swapTx)
	if err != nil {
		e.logger.Println("Failed to send transaction", err)
		return
	}

	e.logger.Println("Trade executed successfully",
		"signature", signature,
		"fromToken", tx.FromToken,
		"toToken", tx.ToToken,
		"amount", amount)
}
