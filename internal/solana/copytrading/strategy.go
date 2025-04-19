package copytrading

import (
	"github.com/pararti/solana-botyara/internal/config"
	"github.com/pararti/solana-botyara/internal/solana"
	"log"
)

// Strategy represents a copytrading strategy
type Strategy struct {
	config *config.Config
	logger *log.Logger

	// Strategy parameters
	maxSlippage  float64
	minTradeSize float64
	maxTradeSize float64
	tradeDelay   int

	// Token lists
	tokenWhitelist map[string]bool
	tokenBlacklist map[string]bool
}

// NewStrategy creates a new copytrading strategy
func NewStrategy(cfg config.Config, logger *log.Logger) *Strategy {
	// Convert token lists to maps for faster lookups
	whitelist := make(map[string]bool)
	for _, token := range cfg.Strategy.TokenWhitelist {
		whitelist[token] = true
	}

	blacklist := make(map[string]bool)
	for _, token := range cfg.Strategy.TokenBlacklist {
		blacklist[token] = true
	}

	return &Strategy{
		logger:         logger,
		maxSlippage:    cfg.Strategy.MaxSlippage,
		minTradeSize:   cfg.Strategy.MinTradeSize,
		maxTradeSize:   cfg.Strategy.MaxTradeSize,
		tradeDelay:     cfg.Strategy.TradeDelay,
		tokenWhitelist: whitelist,
		tokenBlacklist: blacklist,
	}
}

// ShouldCopyTrade determines if a transaction should be copied
func (s *Strategy) ShouldCopyTrade(tx *solana.Transaction) bool {
	// Skip if not a swap
	if !tx.IsSwap {
		return false
	}

	// Check token whitelist (if empty, allow all tokens)
	if len(s.tokenWhitelist) > 0 {
		if !s.tokenWhitelist[tx.FromToken] || !s.tokenWhitelist[tx.ToToken] {
			s.logger.Println("Transaction skipped: token not in whitelist",
				tx.FromToken,
				tx.ToToken)
			return false
		}
	}

	// Check token blacklist
	if s.tokenBlacklist[tx.FromToken] || s.tokenBlacklist[tx.ToToken] {
		s.logger.Println("Transaction skipped: token in blacklist",
			tx.FromToken,
			tx.ToToken)
		return false
	}

	// Check trade size
	if float64(tx.FromAmount) < s.minTradeSize {
		s.logger.Println("Transaction skipped: trade size too small",
			"amount", tx.FromAmount,
			"minTradeSize", s.minTradeSize)
		return false
	}

	if float64(tx.FromAmount) > s.maxTradeSize {
		s.logger.Println("Transaction skipped: trade size too large",
			tx.FromAmount,
			"maxTradeSize", s.maxTradeSize)
		return false
	}

	return true
}

// AdjustTradeAmount adjusts the trade amount based on the strategy
func (s *Strategy) AdjustTradeAmount(originalAmount uint64) uint64 {
	// Implement your logic to adjust the trade amount
	// For example, you might want to trade a percentage of the original amount

	// For now, we'll just copy the exact amount
	return originalAmount
}
