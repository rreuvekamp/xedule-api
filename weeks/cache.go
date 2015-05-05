package weeks

import (
	"time"

	"github.com/rreuvekamp/xedule-api/attendee"
)

type weeksRequest struct {
	ch     chan weeksResponse
	maxAge time.Duration
}

type weeksResponse struct {
	found bool
	wks   weeks
	time  time.Time
}

type cache struct {
	wks  weeks
	time time.Time
}

var defReqMaxAge = time.Minute * 30

var cacheWeeks cache
var att attendee.Attendee

var chWksReq = make(chan weeksRequest, 1)           // Request weeks in cache
var chWksSet = make(chan cache, 1)                  // Set weeks in cache
var chAttReq = make(chan chan attendee.Attendee, 1) // Request attendee for Get

func Run() {
	// Get the first attendee in the database.
	// Used for fetching the weeks list in Get.
	a, err := attendee.FetchS([]string{"id", "name", "type", "lid"},
		"ORDER BY ID LIMIT 1")
	if err == nil || len(a) > 0 {
		att = a[0]
	}

	for {
		select {
		case r := <-chWksReq: // Request cached weeks
			var re weeksResponse
			if time.Since(cacheWeeks.time).Seconds() > r.maxAge.Seconds() {
				r.ch <- re
				continue
			}
			re.found = true
			re.wks = cacheWeeks.wks
			re.time = cacheWeeks.time
			r.ch <- re
		case cw := <-chWksSet:
			cacheWeeks = cw
		case ch := <-chAttReq: // Request attendee used for fetching the weeks list.
			ch <- att
		}
	}
}
