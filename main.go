// MIT License
//
// Copyright (c) 2024-2025 Johnathan A. Hartsfield (George A. Costanza esq.)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package loadb

// This program acts a model load balancer, accepting incoming "requests"
// at a rate moderated using the channel provided via the *time.Ticker type
// (<-&t.T{}.C) and responding to those requests using the same approach.
// As requests come in, we attempt to balance the number of open connections
// (determined by the request rate and response time), over a number of servers
// in different (metaphorical) regions, starting with the next "closest"
// server.
//
// To run the program we use a command like this:
//
// ----------------------------------------------------------------------------
//        $ ./lb --regions=14 --req_rate=13 --res_rate=200 --max_conns=9
// ----------------------------------------------------------------------------
//
//          --regions  ->  the number of regions (servers) to simulate.
//         --req_rate  ->  the number of request to send per millisecond.
//         --res_rate  ->  the number of milliseconds it takes to handle each
//                         request.
//        --max_conns  ->  the maximum number of connection a server can have.
//
// ----------------------------------------------------------------------------
//
// As these parameters are adjusted, one may notice that some configurations
// overload the cluster near instantantly, while some may do so more gradually,
// and, if one is lucky (I don't need luck) one may find sustainable
// configurations, capable of maintaining indefinite integrity.
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
