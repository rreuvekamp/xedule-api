package wsched

import (
	"fmt"
	"strconv"
	"time"
)

type cacheRequest struct {
	ch     chan cacheResponse
	aid    int
	year   int
	week   int
	maxAge time.Duration
}

type cacheResponse struct {
	found bool
	w     WeekSchedule
	time  time.Time
}

type cacheAdd struct {
	w    WeekSchedule
	aid  int
	year int
	week int
	time time.Time
}

type cacheStore struct {
	w    WeekSchedule
	time time.Time
}

var cleanMaxAge = time.Minute * 10  // Max Age used by ticker cleaner.
var defReqMaxAge = time.Minute * 10 // Max Age that can be used with a cache request.

var wscheds = make(map[int]cacheStore)

var chWkAdd = make(chan cacheAdd, 1)     // Add
var chWkReq = make(chan cacheRequest, 1) // Request

// RunCache takes care of cache requests and cleaning up the cache.
func RunCache() {
	// Remove outdated cache items every tick of clean.
	clean := time.NewTicker(10 * time.Minute)

	for {
		select {
		case r := <-chWkReq:
			id := wkId(r.aid, r.year, r.week)
			wk, ok := wscheds[id]

			var re cacheResponse

			if !ok || time.Since(time.Now()).Seconds() > r.maxAge.Seconds() {
				r.ch <- re
				delete(wscheds, id)
				continue
			}

			re.found = true
			re.w = wk.w
			re.time = wk.time
			r.ch <- re
		case r := <-chWkAdd:
			wscheds[wkId(r.aid, r.year, r.week)] = cacheStore{w: r.w, time: r.time}
		case <-clean.C:
			// Check for outdated cache items, and remove them from the map.
			var removes []int
			for id, wk := range wscheds {
				if time.Since(wk.time).Seconds() > cleanMaxAge.Seconds() {
					delete(wscheds, id)
					removes = append(removes, id)
				}
			}
			if len(removes) > 0 {
				fmt.Println("WeekSchedule cache cleaned:", len(removes), removes)
			}
		}
	}
}

func wkId(aid, year, week int) int {
	id, _ := strconv.Atoi(strconv.Itoa(aid) + strconv.Itoa(year-2000) + strconv.Itoa(week))
	return id
}
