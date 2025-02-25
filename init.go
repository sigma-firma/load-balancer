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

import (
	"flag"
	"strconv"
	"time"
)

// region could be a server, or a request to a server.
type region struct {
	LocationID  string
	NextClosest *region
}

// regions could be requests from a region, or a region hosting a
// server which handles incoming requests.
var regions []*region

// activeRegions are regional servers which have active connections.
var activeRegions = make(map[*region]chan *region)

var totalConns int
var count int
var request_rate int
var response_time int
var maxConns int
var start time.Time = time.Now()

// init initializes the requestRegions ([]r), and the r[x].NextClosest property
// which indicates the next geographically close region.
func init() {
	flag.IntVar(&count, "regions", 1234, "The number of clusters to create")
	flag.IntVar(&request_rate, "req_rate", 1234, "Sends a request once every (x) milliseconds")
	flag.IntVar(&response_time, "res_rate", 1234, "Sends a response once every (x) milliseconds")
	flag.IntVar(&maxConns, "max_conns", 1234, "Max connections per node")
	flag.Parse()

	regions = mkRegions(count)
	for i, r := range regions {
		if i+1 == len(regions) {
			r.NextClosest = regions[0]
			return
		}
		r.NextClosest = regions[i+1]
	}
}

func mkRegions(l int) (rs []*region) {
	for i := 1; i <= l; i++ {
		rs = append(rs, &region{"cluster-" + strconv.Itoa(i+100), &region{}})
	}
	return
}
