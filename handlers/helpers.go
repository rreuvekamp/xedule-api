package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type errStr struct {
	Error string
}

// writeJSON writes the given struct (v) on the given http ResponseWriter in JSON format.
// A time object is used for determining last-modified (headers).
func writeJSON(w http.ResponseWriter, r *http.Request, v interface{}, tm time.Time) error {
	indent := ""
	if repeat, err := strconv.Atoi(r.FormValue("indent")); err == nil && repeat > 0 {
		indent = strings.Repeat(" ", repeat)
	}

	var data []byte
	var err error
	if indent == "" {
		data, err = json.Marshal(v)
	} else {
		data, err = json.MarshalIndent(v, "", indent)
	}
	if err != nil {
		return err
	}

	// Check if time is not empty
	if tm != (time.Time{}) {
		if t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since")); err == nil &&
			tm.Unix() <= t.Unix() {

			w.WriteHeader(http.StatusNotModified)
			return errors.New("not modified")
		}
		w.Header().Set("Last-Modified", tm.Format(http.TimeFormat))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err = w.Write(data)
	return err
}
