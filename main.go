package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "expvar"

	"github.com/rreuvekamp/xedule-api/attendee"
	"github.com/rreuvekamp/xedule-api/handlers"
	"github.com/rreuvekamp/xedule-api/lastupdate"
	"github.com/rreuvekamp/xedule-api/misc"
	"github.com/rreuvekamp/xedule-api/weeks"
	"github.com/rreuvekamp/xedule-api/weekschedule"
)

/*
To do:
+/schedule.json&aid=14307&noattinfo=true
+WeekSchedule.BaseUts // UnixTimeStamp of the first second in the week
+event.Start, event.End <- No Uts but amount of seconds since start of day.
+attendee.json?lid=34
+Fixed time zone stuff
+lastupdate.json?year=2015&week=23
+?nocache=true // Does not look for cache, but does update it of course.
 lastupdate and weeks have got their own attendee in cache, can't they share it?
*/

func main() {
	err := misc.LoadConfig(misc.CfgFilename)
	if err != nil {
		os.Exit(1)
	}

	// Don't exit program on error. Without a database this application can still
	// preform some tasks (WeekSchedule without attendee types).
	misc.ConnectDb()

	lid := flag.Int("update-attendees", 0,
		"Location Id of which the attendees should be fetched and updated in database.")
	flag.Parse()
	if *lid > 0 {
		err := attendee.Update(*lid)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		os.Exit(0)
		return
	}

	go wsched.RunCache()
	go weeks.Run()
	go lastupdate.Run()

	http.HandleFunc("/schedule.json", handlers.WSched)
	http.HandleFunc("/weeks.json", handlers.Weeks)
	http.HandleFunc("/attendee.json", handlers.Attendee)
	http.HandleFunc("/lastupdate.json", handlers.LastUpdate)
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("https://github.com/rreuvekamp/xedule-api"))
	})

	log.Println("Started")

	log.Println(http.ListenAndServe(misc.Cfg().Http.Addr, nil))
}
