# Solana Copytrading Bot

This is a Golang-based copytrading bot for the Solana blockchain. It monitors specified wallets for trading activity and automatically copies their trades to your wallet based on a configurable strategy.

## Features

- Monitor multiple wallets for trading activity
- Copy trades with configurable parameters
- Token whitelist and blacklist
- Adjustable trade sizes and slippage
- Logging and error handling

## Requirements

- Go 1.21 or higher
- Access to a Solana RPC endpoint
- A Solana wallet with private key

## Installation

1. Clone the repository:


2. Install dependencies: go mod download




3. Configure the bot by editing `config.yaml`:
- Add your wallet's private key
- Add the wallet addresses you want to copy
- Configure your trading strategy

4. Build the bot: go build -o botyara ./cmd


## Configuration

The bot is configured using the `config.yaml` file:

- `solana`: Solana network configuration
- `endpoint`: RPC endpoint URL
- `network`: Network name (mainnet, testnet, devnet)

- `wallet`: Your wallet configuration
- `private_key`: Your wallet's private key

- `monitoring`: Wallet monitoring configuration
- `wallets`: List of wallet addresses to monitor
- `poll_interval`: How often to check for new transactions (in seconds)
- `confirm_level`: Transaction confirmation level

- `strategy`: Trading strategy configuration
- `type`: Strategy type (mirror, filter, etc.)
- `max_slippage`: Maximum allowed slippage (in %)
- `min_trade_size`: Minimum trade size to copy
- `max_trade_size`: Maximum trade size to copy
- `trade_delay`: Delay before executing a trade (in milliseconds)
- `token_whitelist`: List of token addresses to trade
- `token_blacklist`: List of token addresses to avoid

## Security

- Never share your private key
- Consider running the bot on a secure, dedicated server
- Use a separate wallet for the bot with limited funds

## Disclaimer

This bot is provided for educational purposes only. Trading cryptocurrencies involves significant risk. Use at your own risk.