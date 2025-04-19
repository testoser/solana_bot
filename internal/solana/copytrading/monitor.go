package copytrading

import (
	"context"
	"time"

	"github.com/pararti/solana-botyara/internal/solana"
	"go.uber.org/zap"
)

// Monitor monitors wallets for transactions
type Monitor struct {
	client  *solana.Client
	wallets []string
	logger  *zap.Logger

	// Last seen transactions to avoid duplicates
	lastSeen map[string]time.Time
}

// NewMonitor creates a new wallet monitor
func NewMonitor(client *solana.Client, wallets []string, logger *zap.Logger) *Monitor {
	return &Monitor{
		client:   client,
		wallets:  wallets,
		logger:   logger,
		lastSeen: make(map[string]time.Time),
	}
}

// Start starts monitoring wallets for transactions
func (m *Monitor) Start(ctx context.Context, txCh chan<- *solana.Transaction) {
	m.logger.Info("Starting wallet monitoring", zap.Strings("wallets", m.wallets))

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.checkWallets(ctx, txCh)
		case <-ctx.Done():
			m.logger.Info("Stopping wallet monitoring")
			return
		}
	}
}

// checkWallets checks all monitored wallets for new transactions
func (m *Monitor) checkWallets(ctx context.Context, txCh chan<- *solana.Transaction) {
	for _, wallet := range m.wallets {
		go func(address string) {
			transactions, err := m.client.GetTransactions(ctx, address)
			if err != nil {
				m.logger.Error("Failed to get transactions", zap.Error(err), zap.String("wallet", address))
				return
			}

			for _, tx := range transactions {
				// Skip already seen transactions
				if lastSeen, ok := m.lastSeen[tx.Signature]; ok {
					// If we've seen this transaction in the last hour, skip it
					if time.Since(lastSeen) < time.Hour {
						continue
					}
				}

				// Mark as seen
				m.lastSeen[tx.Signature] = time.Now()

				// Send to channel for processing
				select {
				case txCh <- tx:
					m.logger.Debug("New transaction detected",
						zap.String("wallet", address),
						zap.String("signature", tx.Signature))
				case <-ctx.Done():
					return
				}
			}
		}(wallet)
	}
}
