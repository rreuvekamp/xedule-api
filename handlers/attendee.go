package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rreuvekamp/xedule-api/attendee"
)

func Attendee(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var idsStr string

	for k, vs := range r.Form {
		switch k {
		case "aid", "ids":
			for _, v := range vs {
				split := strings.Split(v, ",")

			loopSplit:
				for _, idStr := range split {
					_, err := strconv.Atoi(idStr)
					if err != nil {
						continue loopSplit
					}

					if idsStr != "" {
						idsStr += ", "
					}
					idsStr += idStr
				}
			}
		}
	}

	lid, _ := strconv.Atoi(r.FormValue("lid"))

	var sqlStr string
	switch true {
	case idsStr != "":
		sqlStr = "WHERE id IN (" + idsStr + ")"
	case lid > 0:
		sqlStr = "WHERE lid = " + strconv.Itoa(lid)
	default:
		writeJSON(w, r, errStr{Error: "invalid aid/ids (attendee id) and lid (location id)"}, time.Time{})
		return
	}

	atts, err := attendee.FetchS([]string{"id", "name", "type"}, sqlStr)

	if err != nil {
		writeJSON(w, r, errStr{Error: "error fetching attendees"}, time.Time{})
		return
	}

	writeJSON(w, r, atts, time.Time{})
}
