package main

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

// reSelect is a recursive function which attempts to send a request to a
// reqion, defaulting to the region.NextClosest when the initial reqion has
// reached capacity
func reSelect(l *reqion) {
	select {
	case activeRegions[l] <- &reqion{mkID(15), l}:
	default:
		reSelect(l.NextClosest)
	}
}

// mkID returns a unique identifier of length idLen
func mkID(idLen int) (id string) {
	symbols := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	for i := 0; i <= idLen; i++ {
		id = id + string(symbols[rand.Intn(len(symbols))])
	}
	return id
}

// perMs returns a *time.Ticker, ticking "per" millisecond
func perMs(per time.Duration) *time.Ticker {
	return time.NewTicker(per * time.Millisecond)
}

// display is used for visualizing the load balancing process via terminal
// output
func display(srvs map[*reqion]chan *reqion) {
	for i := 0; i <= 100; i++ {
		// clear the terminal
		fmt.Println()
	}
	for _, loc := range requestRegions {
		vlen := ""
		for i := 0; i <= len(srvs[loc]); i++ {
			vlen = vlen + "="
		}
		fmt.Println(loc.LocationID, len(srvs[loc]), "\t", vlen)
	}
}
