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
	"github.com/rreuvekamp/xedule-api/misc"
	"github.com/rreuvekamp/xedule-api/weeks"
	"github.com/rreuvekamp/xedule-api/weekschedule"
)

/*
To do:
 Have attendees be in memory instead of database.
 Log HTTP requests
 /attendee.json?aid=14327,14309
+/schedule.json also giving list of attendee ids
*/

func main() {
	err := misc.LoadConfig(misc.CfgFilename)
	if err != nil {
		os.Exit(1)
	}

	// Don't exit program. Without a database this application can still
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

	http.HandleFunc("/schedule.json", handlers.WSched)
	http.HandleFunc("/weeks.json", handlers.Weeks)
	http.HandleFunc("/attendee.json", handlers.Attendee)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("https://github.com/rreuvekamp/xedule-api"))
	})

	log.Println("Started")

	log.Println(http.ListenAndServe(misc.Cfg().Http.Addr, nil))
}
