package handlers

import (
	"net/http"

	"github.com/rreuvekamp/xedule-api/weeks"
)

func Weeks(w http.ResponseWriter, r *http.Request) {
	wks, tm, _ := weeks.Get()
	writeJSON(w, r, wks, tm)
}
