package main

import (
	"fmt"

	"github.com/rreuvekamp/xedule-api/misc"
	"github.com/rreuvekamp/xedule-api/types/attendee"
)

/*
To do:
+attendee.Update for fetching all attendees for a given location.
 Rewrite wsched.Fetch without external ICS parser.
*/

func main() {
	misc.ConnectDb()
	//fmt.Println(attendee.FetchS([]string{}, ""))
	//fmt.Println(wsched.Fetch(14327, 2015, 18))
	fmt.Println(attendee.Update(34))
}
