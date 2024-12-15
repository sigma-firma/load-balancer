package main

// reqion could be a server, or a request to a server.
type reqion struct {
	LocationID  string
	NextClosest *reqion
}

// requestRegions could be requests from a region, or a region which hosts a
// servers which handle the requests.
var requestRegions []*reqion = []*reqion{
	{"US", &reqion{}}, {"EU", &reqion{}}, {"CH", &reqion{}}, {"PH", &reqion{}},
	{"JP", &reqion{}}, {"AR", &reqion{}}, {"CA", &reqion{}}, {"IS", &reqion{}},
	{"GB", &reqion{}}, {"SA", &reqion{}}, {"MC", &reqion{}}, {"RU", &reqion{}},
	{"AU", &reqion{}}, {"NZ", &reqion{}}, {"AN", &reqion{}},
}

// activeRegions are regional servers which have active connections.
var activeRegions = make(map[*reqion]chan *reqion)

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
