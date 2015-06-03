package handlers

import (
	"net/http"

	"github.com/rreuvekamp/xedule-api/weeks"
)

func Weeks(w http.ResponseWriter, r *http.Request) {

	cache := true

	// NoCache only has effect if the remoteaddr is whitelisted in the config file.
	if r.FormValue("nocache") != "" && checkCacheWhitelist(ip(r)) {
		cache = false
	}

	wks, tm, _ := weeks.Get(cache)

	writeJSON(w, r, wks, tm)
}
