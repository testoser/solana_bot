package solana

import (
	"context"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// Client represents a Solana RPC client
type Client struct {
	rpcClient *rpc.Client
	logger    *log.Logger
}

// NewClient creates a new Solana client
func NewClient(endpoint string, logger *log.Logger) (*Client, error) {
	rpcClient := rpc.New(endpoint)

	// Test connection
	_, err := rpcClient.GetVersion(context.Background())
	if err != nil {
		return nil, err
	}

	return &Client{
		rpcClient: rpcClient,
		logger:    logger,
	}, nil
}

// GetBalance gets the balance of a wallet
func (c *Client) GetBalance(ctx context.Context, address string) (uint64, error) {
	pubKey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return 0, err
	}

	balance, err := c.rpcClient.GetBalance(
		ctx,
		pubKey,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return 0, err
	}

	return balance.Value, nil
}

// GetTransactions gets recent transactions for a wallet
func (c *Client) GetTransactions(ctx context.Context, address string) ([]*Transaction, error) {
	pubKey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return nil, err
	}

	signatures, err := c.rpcClient.GetSignaturesForAddress(
		ctx,
		pubKey,
	)
	if err != nil {
		return nil, err
	}

	var transactions []*Transaction
	for _, sig := range signatures {
		tx, err := c.GetTransaction(ctx, sig.Signature.String())
		if err != nil {
			c.logger.Println("Failed to get transaction", sig.Signature.String())
			continue
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// GetTransaction gets a transaction by signature
func (c *Client) GetTransaction(ctx context.Context, signature string) (*Transaction, error) {
	sig, err := solana.SignatureFromBase58(signature)
	if err != nil {
		return nil, err
	}

	tx, err := c.rpcClient.GetTransaction(
		ctx,
		sig,
		&rpc.GetTransactionOpts{
			Encoding:   solana.EncodingBase64,
			Commitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		return nil, err
	}

	// Process and return transaction data
	bt := tx.BlockTime.Time().Unix()
	return &Transaction{
		Signature: signature,
		BlockTime: &bt,
		Slot:      tx.Slot,
		// Process transaction data as needed
	}, nil
}

// SendTransaction sends a transaction to the network
func (c *Client) SendTransaction(ctx context.Context, tx *solana.Transaction) (string, error) {
	sig, err := c.rpcClient.SendTransactionWithOpts(
		ctx,
		tx,
		rpc.TransactionOpts{
			SkipPreflight:       false,
			PreflightCommitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		return "", err
	}

	return sig.String(), nil
}
