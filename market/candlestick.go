package market

import (
	"fmt"
)

// SetOpen sets the open price to the supplied value, result is returned as a
// new candlestick.
func (c Candlestick) SetOpen(open float64) Candlestick {
	c.Open = open
	return c
}

// SetClose sets the close price to the supplied value, result is returned as a
// new candlestick.
func (c Candlestick) SetClose(close float64) Candlestick {
	c.Close = close
	return c
}

// SetHigh sets the high price to the supplied value, result is returned as a
// new candlestick.
func (c Candlestick) SetHigh(high float64) Candlestick {
	c.High = high
	return c
}

// SetLow sets the low price to the supplied value, result is returned as a new
// candlestick.
func (c Candlestick) SetLow(low float64) Candlestick {
	c.Low = low
	return c
}

// SetVolume sets the volume to the supplied value, result is returned as a new
// candlestick.
func (c Candlestick) SetVolume(volume int) Candlestick {
	c.Volume = volume
	return c
}

// SetOpensHour sets whether this candlestick opens a new hour.
func (c Candlestick) SetOpensHour(opensHour bool) Candlestick {
	c.OpensHour = opensHour
	return c
}

// SetOpensDay sets whether this candlestick opens a new day.
func (c Candlestick) SetOpensDay(opensDay bool) Candlestick {
	c.OpensDay = opensDay
	return c
}

// Add adds the supplied values to this candlestick, result is returned as a new
// candlestick.
func (c Candlestick) Add(open, close, high, low float64,
	volume int) Candlestick {
	c.Open += open
	c.Close += close
	c.High += high
	c.Low += low
	c.Volume += volume
	return c
}

// AddCandlestick adds the values of another candlestick to this candlestick,
// result is returned as a new candlestick.
func (c Candlestick) AddCandlestick(other Candlestick) Candlestick {
	c.Open += other.Open
	c.Close += other.Close
	c.High += other.High
	c.Low += other.Low
	c.Volume += other.Volume
	return c
}

// Subtract subtracts the supplied values from this candlestick, result is
// returned as a new candlestick.
func (c Candlestick) Subtract(open, close, high, low float64,
	volume int) Candlestick {
	c.Open -= open
	c.Close -= close
	c.High -= high
	c.Low -= low
	c.Volume -= volume
	return c
}

// SubtractCandlestick subtracts the values of another candlestick from this
// candlestick, result is returned as a new candlestick.
func (c Candlestick) SubtractCandlestick(other Candlestick) Candlestick {
	c.Open -= other.Open
	c.Close -= other.Close
	c.High -= other.High
	c.Low -= other.Low
	c.Volume -= other.Volume
	return c
}

// Multiply multiplies the supplied values to this candlestick, result is
// returned as a new candlestick.
func (c Candlestick) Multiply(open, close, high, low float64,
	volume int) Candlestick {
	c.Open *= open
	c.Close *= close
	c.High *= high
	c.Low *= low
	c.Volume *= volume
	return c
}

// MultiplyCandlestick multiplies the values of another candlestick to this
// candlestick, result is returned as a new candlestick.
func (c Candlestick) MultiplyCandlestick(other Candlestick) Candlestick {
	c.Open *= other.Open
	c.Close *= other.Close
	c.High *= other.High
	c.Low *= other.Low
	c.Volume *= other.Volume
	return c
}

// Divide divides this candlestick by the supplied values, result is returned as
// a new candlestick.
func (c Candlestick) Divide(open, close, high, low float64,
	volume int) Candlestick {
	c.Open /= open
	c.Close /= close
	c.High /= high
	c.Low /= low
	c.Volume /= volume
	return c
}

// DivideCandlestick divides this candlestick by the values of another
// candlestick, result is returned as a new candlestick.
func (c Candlestick) DivideCandlestick(other Candlestick) Candlestick {
	c.Open /= other.Open
	c.Close /= other.Close
	c.High /= other.High
	c.Low /= other.Low
	c.Volume /= other.Volume
	return c
}

// String returns a string representation of this candlestick.
func (c Candlestick) String() string {
	return fmt.Sprintf(
		"Exchange: %s Ticker: %s Created At: %v, Open: %.2f, Close: %.2f, High: %.2f Low: %.2f Volume: %d",
		c.Exchange, c.Ticker, c.CreatedAt, c.Open, c.Close, c.High, c.Low,
		c.Volume)
}
