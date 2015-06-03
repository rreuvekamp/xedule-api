package lastupdate

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rreuvekamp/xedule-api/attendee"
)

type cacheReq struct {
	ch   chan cacheResp
	year int
	week int
}

type cacheResp struct {
	found bool
	tmLu  time.Time
	tm    time.Time // Cache time
}

type cacheAdd struct {
	tmLu time.Time
	year int
	week int
	tm   time.Time // Cache time
}

type cacheItem struct {
	tmLu time.Time
	year int
	week int
	tm   time.Time // Cache time
}

var lastupdates = make(map[int]cacheItem)
var att attendee.Attendee

var cleanMaxAge = 10 * time.Minute

var chLuReq = make(chan cacheReq, 1) // Request (/check for) cache
var chLuAdd = make(chan cacheAdd, 1) // Add cache
var chAttReq = make(chan chan attendee.Attendee)

// Run inits cache and handles cache requests.
func Run() {
	// Get the first attendee in the database.
	// Used for fetching the weeks list in Get.
	a, err := attendee.FetchS([]string{"id", "name", "type", "lid"},
		"ORDER BY ID LIMIT 1")
	if err == nil || len(a) > 0 {
		att = a[0]
	}

	clean := time.NewTicker(10 * time.Minute)
	for {
		select {
		case r := <-chLuReq: // Request (/lookup) cache
			ci, ok := lastupdates[luId(r.year, r.week)]

			var re cacheResp

			if !ok {
				r.ch <- re
				continue
			}

			re.found = true
			re.tmLu = ci.tmLu
			re.tm = ci.tm // Cache time

			r.ch <- re
		case r := <-chLuAdd: // Add/update cache
			lastupdates[luId(r.year, r.week)] = cacheItem{tmLu: r.tmLu, year: r.year, week: r.week, tm: r.tm}
		case ch := <-chAttReq:
			ch <- att
		case <-clean.C: // Periodicly clean up cache
			// Check for outdated cache items, and remove them from the map.
			var removes []int
			for id, ci := range lastupdates {
				if time.Since(ci.tm).Seconds() > cleanMaxAge.Seconds() {
					delete(lastupdates, id)
					removes = append(removes, id)
				}
			}
			if len(removes) > 0 {
				fmt.Println("LastUpdates cache cleaned:", len(removes), removes)
			}
		}
	}
}

// luId makes an id used in the cache map by year/week.
func luId(year, week int) int {
	id, _ := strconv.Atoi(strconv.Itoa(year) + strconv.Itoa(week))
	return id
}
