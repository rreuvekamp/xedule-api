package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/rreuvekamp/xedule-api/weekschedule"
)

// WSched HTTP handler gives the WeekSchedule for the given attendee id (aid), year and week.
// Defaults to current year and week (next week if current weekday >= saterday).
// Format changes by setting 'legacy' form value.
func WSched(w http.ResponseWriter, r *http.Request) {
	var aid, year, week int
	var err error

	tm := time.Now()
	// Next week if saterday.
	if tm.Weekday() == 6 {
		tm = tm.AddDate(0, 0, 3)
	}
	cYear, cWeek := tm.ISOWeek()

	if aid, err = strconv.Atoi(r.FormValue("aid")); err != nil || aid <= 0 {
		writeJSON(w, r, errStr{Error: "invalid aid (attendee id)"}, time.Time{})
		return
	}

	// If given year is <= 2010 then obviously not a serieus request.
	if year, err = strconv.Atoi(r.FormValue("year")); err != nil || year <= 2010 {
		year = cYear
	}

	if week, err = strconv.Atoi(r.FormValue("week")); err != nil || week <= 0 {
		week = cWeek
	}

	wk, tm, _ := wsched.Get(aid, year, week)

	if r.FormValue("legacy") != "" {
		writeJSON(w, r, wk.Legacy(), tm)
		return
	}
	writeJSON(w, r, wk, tm)
}
