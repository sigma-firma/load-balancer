package main

import (
	"time"

	"golang.org/x/exp/rand"
)

func main() {
	var (
		// incoming requests per millisecond
		requests *time.Ticker = perMs(time.Duration(request_rate))
		// millisecond avg response time
		resTime *time.Ticker  = perMs(time.Duration(response_time))
		quit    chan struct{} = make(chan struct{})
	)
	for {
		select {
		case <-requests.C: // on each request
			// pick a random region
			l := regions[rand.Intn(len(regions))]
			// if the region has no connections, we open a channel
			// which handles requests for that region, creating a
			// new activeRegion
			if _, ok := activeRegions[l]; !ok {
				activeRegions[l] = make(chan *region, maxConns)
			}
			// We send the request to the region/server, if the
			// server connections pool is at capacity, reSelect()
			// recursively tries the region.NextClosest region.
			// see: reSelect()
			reSelect(l)
			display(activeRegions) // terminal output
		case <-resTime.C: // request response
			for l, r := range activeRegions {
				<-r
				totalConns = totalConns - 1
				if len(r) == 0 {
					// if all requests have been handled we
					// delete the region from activeRegions
					delete(activeRegions, l)
				}
			}
		case <-quit:
			requests.Stop()
			resTime.Stop()
			return
		}
	}
}
