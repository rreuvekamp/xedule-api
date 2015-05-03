package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/rreuvekamp/xedule-api/handlers"
	"github.com/rreuvekamp/xedule-api/misc"
	"github.com/rreuvekamp/xedule-api/types/attendee"
	"github.com/rreuvekamp/xedule-api/types/weekschedule"
)

/*
To do:
 Have attendees be in memory instead of database.
+Cache of WeekSchedule clean up every 15 minutes.
 Log HTTP requests
+Readme markdown
+Config file
*/

func main() {
	err := misc.LoadConfig(misc.CfgFilename)
	if err != nil {
		os.Exit(1)
	}

	// Don't exit program. Without a database this application can still preform some tasks (WeekSchedule without attendee types).
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

	http.HandleFunc("/schedule.json", handlers.WSched)

	fmt.Println(http.ListenAndServe(misc.Cfg().Http.Addr, nil))
}
