package main

import (
	"golang.org/x/exp/rand"
)

func main() {
	requests := perMs(13) // incoming requests per millisecond
	resTime := perMs(200) // millisecond avg response time
	quit := make(chan struct{})
	for {
		select {
		case <-requests.C: // create each request
			// pick a random region
			l := requestRegions[rand.Intn(len(requestRegions))]
			// if the region has no connections, we open a channel
			// which handles requests for that region, creating a
			// new activeRegion
			if _, ok := activeRegions[l]; !ok {
				activeRegions[l] = make(chan *region, 9)
			}
			// We send the request to the region/server, if the
			// server connections pool is at capacity, reSelect()
			// recursively tries the region.NextClosest reqion.
			// see: reSelect()
			reSelect(l)
			// terminal output
			display(activeRegions)
		case <-resTime.C: // handle each request
			for l, r := range activeRegions {
				<-r
				if len(r) == 0 {
					// if all requests have been handled we
					// delete the reqion from activeRegions
					delete(activeRegions, l)
				}
			}
		case <-quit:
			requests.Stop()
			return
		}
	}

}
