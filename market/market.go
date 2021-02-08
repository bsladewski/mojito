package market

// PlatformKey refers to a specific platform for retrieving market data and
// facilitating trades.
type PlatformKey string

// Define supported platforms.
const (
	PlatformCoinbase PlatformKey = "coinbase"
	PlatformAlpaca   PlatformKey = "alpaca"
)

// ExchangeKey refers to a specific exchange for trading securities.
type ExchangeKey string

// Define supported exchanges.
const (
	ExchangeCoinbase     ExchangeKey = "COINBASE"
	ExchangeIEX          ExchangeKey = "IEX"
	ExchangeNASDAQBX     ExchangeKey = "NASDAQ_BX"
	ExchangeNASDAQPSX    ExchangeKey = "NASDAQ_PSX"
	ExchangeNYSENational ExchangeKey = "NYSE_NATIONAL"
	ExchangeNYSEChicago  ExchangeKey = "NYSE_CHICAGO"
)
