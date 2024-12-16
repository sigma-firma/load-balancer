package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/exp/rand"
)

// reSelect is a recursive function which attempts to send a request to a
// reqion, defaulting to the region.NextClosest when the initial reqion has
// reached capacity
func reSelect(l *region) {
	select {
	case activeRegions[l] <- &region{mkID(15), l}:
		totalConns = totalConns + 1
		if totalConns >= maxConns*count {
			os.Exit(0)
		}
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
func display(srvs map[*region]chan *region) {
	// clear the terminal
	for i := 0; i <= 10; i++ {
		fmt.Print("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
	}

	// print variables
	fmt.Println("->Regions:\t", count)
	fmt.Println("->Connections:\t", totalConns, "/", count*maxConns)
	fmt.Println("->Requests/ms:\t", request_rate)
	fmt.Println("->Responses/ms:\t", response_rate)
	fmt.Println("->Uptime:\t", time.Since(start))
	fmt.Println()

	for _, loc := range regions {
		// visualize connection count/dispersement
		vlen := ""
		for i := 0; i <= len(srvs[loc]); i++ {
			vlen = vlen + "="
		}
		// print ID and connection count
		fmt.Println(loc.LocationID, "|", len(srvs[loc]), "|", vlen)
	}
	fmt.Println()
	fmt.Println()
}
