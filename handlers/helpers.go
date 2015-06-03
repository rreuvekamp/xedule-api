package handlers

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rreuvekamp/xedule-api/misc"
)

type errStr struct {
	Error string `json:"error"`
}

// writeJSON writes the given object (v) on the given http ResponseWriter in JSON format.
// The time given should be the time the object was last modified, for setting the last-modified
// headers and sending StatusNotModified if the client has the correct version already.
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
		t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since"))
		if err == nil && tm.Unix() <= t.Unix() {
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

// ip returns the (correct) Ip address for the given http.Request.
func ip(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	// Proxy stuff
	if ip == "127.0.0.1" {
		ip = r.Header.Get("X-FORWARDED-FOR")
	}
	return ip
}

// checkCacheWhitelist returns a boolean which is true if the given IP address is
// whitelisted to make API requests of which cache is not looked up.
func checkCacheWhitelist(addr string) bool {
	// NoCache only has effect if the remoteaddr is whitelisted in the config file.
	for _, a := range misc.Cfg().Api.NoCacheIpAllow {
		if a == addr {
			return true
		}
	}
	return false
}
