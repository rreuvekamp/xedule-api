package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/rreuvekamp/xedule-api/lastupdate"
)

type pageLu struct {
	Year int   `json:"year"`
	Week int   `json:"week"`
	Uts  int64 `json:"uts"`
}

func LastUpdate(w http.ResponseWriter, r *http.Request) {
	yr, wk := time.Now().ISOWeek()

	if year, err := strconv.Atoi(r.FormValue("year")); err == nil {
		yr = year
	}

	if week, err := strconv.Atoi(r.FormValue("week")); err == nil {
		wk = week
	}

	cache := true

	// NoCache only has effect if the remoteaddr is whitelisted in the config file.
	if r.FormValue("nocache") != "" && checkCacheWhitelist(ip(r)) {
		cache = false
	}

	tm, tmCache, _ := lastupdate.Get(yr, wk, cache)

	if tm == (time.Time{}) {
		p := errStr{
			Error: "time for year/week not found",
		}
		writeJSON(w, r, p, time.Time{})
		return
	}

	p := pageLu{
		Year: yr,
		Week: wk,
		Uts:  tm.Unix(),
	}

	writeJSON(w, r, p, tmCache)
}
