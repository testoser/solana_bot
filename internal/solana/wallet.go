package solana

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/gagliardetto/solana-go"
)

// Wallet represents a Solana wallet
type Wallet struct {
	keypair solana.PrivateKey
	pubKey  solana.PublicKey
	logger  *log.Logger
}

// NewWalletFromPrivateKey creates a new wallet from a private key
func NewWalletFromPrivateKey(privateKeyBase58 string, logger *log.Logger) (*Wallet, error) {
	// Check if the key is in base58 or base64 format
	var keypair solana.PrivateKey
	var err error

	// Try base58 first
	keypair, err = solana.PrivateKeyFromBase58(privateKeyBase58)
	if err != nil {
		// Try base64
		keyBytes, err := base64.StdEncoding.DecodeString(privateKeyBase58)
		if err != nil {
			return nil, err
		}
		keypair = solana.PrivateKey(keyBytes)
	}

	pubKey := keypair.PublicKey()

	logger.Println("Wallet initialized", pubKey.String())

	return &Wallet{
		keypair: keypair,
		pubKey:  pubKey,
		logger:  logger,
	}, nil
}

// Address returns the wallet address
func (w *Wallet) Address() string {
	return w.pubKey.String()
}

// Sign signs a transaction
func (w *Wallet) Sign(tx *solana.Transaction) error {
	_, err := tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(w.pubKey) {
				return &w.keypair
			}
			return nil
		},
	)
	return err
}

// CreateSwapTransaction creates a transaction to swap tokens
func (w *Wallet) CreateSwapTransaction(
	ctx context.Context,
	client *Client,
	fromToken string,
	toToken string,
	amount uint64,
) (*solana.Transaction, error) {
	// This is a simplified example. In a real implementation, you would:
	// 1. Find the best route for the swap (using Jupiter or another DEX aggregator)
	// 2. Build the swap instructions
	// 3. Create and sign the transaction

	// For demonstration purposes, we'll create a placeholder transaction
	recentBlockhash, err := client.rpcClient.GetRecentBlockhash(ctx, "finalized")
	if err != nil {
		return nil, err
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			// Add swap instructions here
		},
		recentBlockhash.Value.Blockhash,
	)

	return tx, err
}
