package copytrading

import (
	"context"
	"log"
	"time"

	"github.com/pararti/solana-botyara/internal/solana"
)

// Monitor monitors wallets for transactions
type Monitor struct {
	client  *solana.Client
	wallets []string
	logger  *log.Logger

	// Last seen transactions to avoid duplicates
	lastSeen map[string]time.Time
}

// NewMonitor creates a new wallet monitor
func NewMonitor(client *solana.Client, wallets []string, logger *log.Logger) *Monitor {
	return &Monitor{
		client:   client,
		wallets:  wallets,
		logger:   logger,
		lastSeen: make(map[string]time.Time),
	}
}

// Start starts monitoring wallets for transactions
func (m *Monitor) Start(ctx context.Context, txCh chan<- *solana.Transaction) {
	m.logger.Println("Starting wallet monitoring", "wallets", m.wallets)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.checkWallets(ctx, txCh)
		case <-ctx.Done():
			m.logger.Println("Stopping wallet monitoring")
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
				m.logger.Println("Failed to get transactions", err, "wallet")
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
					m.logger.Println("New transaction detected",
						"wallet", address,
						"signature", tx.Signature)
				case <-ctx.Done():
					return
				}
			}
		}(wallet)
	}
}
