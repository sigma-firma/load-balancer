package main

// region could be a server, or a request to a server.
type region struct {
	LocationID  string
	NextClosest *region
}

// requestRegions could be requests from a region, or a region which hosts a
// servers which handle the requests.
var requestRegions []*region = []*region{
	{"US", &region{}}, {"EU", &region{}}, {"CH", &region{}}, {"PH", &region{}},
	{"JP", &region{}}, {"AR", &region{}}, {"CA", &region{}}, {"IS", &region{}},
	{"GB", &region{}}, {"SA", &region{}}, {"MC", &region{}}, {"RU", &region{}},
	{"AU", &region{}}, {"NZ", &region{}}, {"AN", &region{}},
}

// activeRegions are regional servers which have active connections.
var activeRegions = make(map[*region]chan *region)

// init initializes the requestRegions ([]r), and the r[x].NextClosest property
// which indicates the next geographically close region.
func init() {
	for i, r := range requestRegions {
		if i+1 < len(requestRegions) {
			r.NextClosest = requestRegions[i+1]
		} else {
			r.NextClosest = requestRegions[0]
		}
	}
}
