// Package leadingedgeticker contains a wrapper around time.Ticker that
// immediately ticks when created.
package leadingedgeticker

import "time"

// A LeadingEdgeTicker holds a time.Ticker that delivers `ticks' of a clock
// at intervals. It will also send a tick when created.
type LeadingEdgeTicker struct {
	C           <-chan time.Time
	innerTicker *time.Ticker
}

// NewTicker returns a new leading edge Ticker containing a channel that
// will send the time with a period specified by the duration argument.
// It adjusts the intervals or drops ticks to make up for slow receivers.
// The duration d must be greater than zero; if not, NewTicker will panic.
// Stop the ticker to release associated resources.
func NewTicker(d time.Duration) LeadingEdgeTicker {
	innerTicker := time.NewTicker(d)

	outputChan := make(chan time.Time)

	go func() {
		outputChan <- time.Now()
		for {
			v, ok := <-innerTicker.C
			if !ok {
				return
			}
			outputChan <- v
		}
	}()

	return LeadingEdgeTicker{
		C:           outputChan,
		innerTicker: innerTicker,
	}
}

// Stop turns off a ticker. After Stop, no more ticks will be sent.
// Stop does not close the channel, to prevent a read from the channel succeeding
// incorrectly.
func (t LeadingEdgeTicker) Stop() {
	t.innerTicker.Stop()
}
