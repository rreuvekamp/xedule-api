package wsched

import (
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

var defMaxAge = time.Minute * 10

var wscheds = make(map[int]cacheStore)

var chWkReq = make(chan cacheRequest, 1)
var chWkAdd = make(chan cacheAdd, 1)

func RunCache() {
	for {
		select {
		case r := <-chWkReq:
			wk, ok := wscheds[wkId(r.aid, r.year, r.week)]

			var re cacheResponse
			if !ok {
				r.ch <- re
				continue
			}

			re.found = true
			re.w = wk.w
			re.time = wk.time
			r.ch <- re
		case r := <-chWkAdd:
			wscheds[wkId(r.aid, r.year, r.week)] = cacheStore{w: r.w, time: r.time}
		}
	}
}

func wkId(aid, year, week int) int {
	id, _ := strconv.Atoi(strconv.Itoa(aid) + strconv.Itoa(year-2000) + strconv.Itoa(week))
	return id
}
