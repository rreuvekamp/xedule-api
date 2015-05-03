package wsched

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/rreuvekamp/xedule-api/types/attendee"
)

// WeekSchedule contains all Days and Events of an attendee for a week.
type WeekSchedule struct {
	Year int           `json:"year"`
	Week int           `json:"week"`
	Days []DaySchedule `json:"days"`
}

// DaySchedule contains all Events of an attendee for a day.
type DaySchedule struct {
	Day    time.Weekday `json:"day"`
	Events []Event      `json:"events"`
}

// Event is a single event for an attendee.
type Event struct {
	Start   int64    `json:"start"`
	End     int64    `json:"end"`
	Desc    string   `json:"desc"`              // Description
	Classes []string `json:"classes,omitempty"` // (Other) classes/attendees
	Facs    []string `json:"facs,omitempty"`    // Facilities
	Staffs  []string `json:"staffs,omitempty"`

	// Used by Fetch
	atts []string

	start time.Time
	end   time.Time
}

const icsTimeLayout = "20060102T150405Z"
const urlWSched = "https://summacollege.xedule.nl/Calendar/iCalendarICS/%d?year=%d&week=%d"

// Get either returns the WeekSchedule for the given aid, year and week from cache
// or if no valid cache, fetches the ICS file, parses it and returns the WeekSchedule from that.
func Get(aid, year, week int) (WeekSchedule, time.Time, error) {

	// Request cache
	ch := make(chan cacheResponse)
	chWkReq <- cacheRequest{
		ch:     ch,
		aid:    aid,
		year:   year,
		week:   week,
		maxAge: defReqMaxAge,
	}

	// Wait for and handle cache response
	c := <-ch
	if c.found {
		return c.w, c.time, nil
	}

	// Check cache
	// Serve cache if not outdated.

	resp, err := http.Get(fmt.Sprintf(urlWSched, aid, year, week))
	if err != nil {
		log.Println("ERROR fetching ICS file of weekschedule:", err, aid, year, week)
		return WeekSchedule{}, time.Time{}, err
	}
	defer resp.Body.Close()

	ics, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ERROR reading from fetched ICS file of weekschedule:", err, aid, year, week)
		return WeekSchedule{}, time.Time{}, err
	}

	// Variables for the parsing
	var cur Event
	var days []DaySchedule
	var atts []string // Slice of attendee names which type should be looked up.

	// The parsing itself
loop:
	for _, l := range strings.Split(string(ics), "\n") {
		var err error
		split := strings.Split(l, ":")
		title := strings.Split(split[0], ";")[0] // Extra split for semi-colon for ATTENDEE
		switch title {
		case "BEGIN":
			if icsIndex(split, 1) != "VEVENT" {
				continue
			}
			// Clean/reset current event.
			cur = Event{}
		case "DTSTART": // DTEND:20150428T090000Z
			cur.start, err = time.Parse(icsTimeLayout, strings.TrimSpace(icsIndex(split, 1)))
			if err != nil {
				log.Println("ERROR parsing start time of ICS: \n", err, split)
				continue loop
			}
			cur.Start = cur.start.Unix()
		case "DTEND": // DTSTART:20150428T073000Z
			cur.end, err = time.Parse(icsTimeLayout, strings.TrimSpace(icsIndex(split, 1)))
			if err != nil {
				log.Println("ERROR parsing end time of ICS: \n", err, split)
				continue loop
			}
			cur.End = cur.end.Unix()
		case "DESCRIPTION": // DESCRIPTION:test
			desc := icsIndex(split, 1)
			if len(desc) != 0 {
				cur.Desc = desc
			}
		case "LOCATION": // LOCATION:BA6.00
			loc := icsIndex(split, 1)
			if len(loc) != 0 {
				cur.Facs = append(cur.Facs, loc)
			}
		case "ATTENDEE": // ATTENDEE;CN=XED:MAILTO:noreply@xedule.nl
			split = strings.Split(icsIndex(split, 0), ";")
			cur.atts = append(cur.atts, strings.TrimPrefix(icsIndex(split, 1), "CN="))
		case "END":
			if icsIndex(split, 1) != "VEVENT" {
				continue
			}

			// Append cur to Day if it exists already, or append it to a newly created Day.
			var success bool
			for i, d := range days {
				if d.Day == cur.start.Weekday() {
					days[i].Events = append(d.Events, cur)
					success = true
				}
			}
			if !success {
				days = append(days, DaySchedule{
					Events: []Event{cur},
					Day:    cur.start.Weekday(),
				})
			}

			atts = append(atts, cur.atts...)
		}
	}

	sort.Sort(DaysByDay(days))

	for di, _ := range days {
		sort.Sort(EventsByStart(days[di].Events))
	}

	w := WeekSchedule{
		Days: days,
		Year: year,
		Week: week,
	}

	if len(atts) > 0 {
		w.findAtts(atts)
	}

	chWkAdd <- cacheAdd{w: w, aid: aid, year: year, week: week, time: time.Now()}

	return w, time.Time{}, nil
}

// findAtts sorts the attendees of the WeekSchedule's events
// as class, staff or facility by looking the names up in the database.
func (w *WeekSchedule) findAtts(names []string) error {
	// Format the names.
	var end string
	for i, n := range names {
		if i > 0 {
			end += ", "
		}
		end += "'" + n + "'"
	}

	// Query the attendees by name.
	atts, err := attendee.FetchS([]string{"name", "type"}, "WHERE name IN ("+end+")")
	if err != nil {
		return err
	}

	// Make a map of atts, for easier lookup by name below.
	attsM := make(map[string]attendee.Attendee)
	for _, a := range atts {
		attsM[a.Name] = a
	}

	for di, d := range w.Days { // Day
		for ei, e := range d.Events { // Event
			for _, ea := range e.atts { // Attendee of the event (EventAttendee)
				// Lookup attendee by name.
				att, ok := attsM[ea]
				if !ok {
					continue
				}

				// Determine which type this attendee is.
				switch att.Type {
				case attendee.Class:
					w.Days[di].Events[ei].Classes = append(w.Days[di].Events[ei].Classes, att.Name)
				case attendee.Staff:
					w.Days[di].Events[ei].Staffs = append(w.Days[di].Events[ei].Staffs, att.Name)
				case attendee.Facil:
					w.Days[di].Events[ei].Facs = append(w.Days[di].Events[ei].Facs, att.Name)

				}
			}
		}
	}

	return nil
}

// icsIndex is a helper function that returns the strings at the given index of the given slice.
func icsIndex(slice []string, index int) string {
	if len(slice) > index {
		return strings.TrimSpace(slice[index])
	}
	return ""
}

// Sorting functions

type DaysByDay []DaySchedule

func (d DaysByDay) Len() int           { return len(d) }
func (d DaysByDay) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DaysByDay) Less(i, j int) bool { return d[i].Day < d[j].Day }

type EventsByStart []Event

func (e EventsByStart) Len() int           { return len(e) }
func (e EventsByStart) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e EventsByStart) Less(i, j int) bool { return e[i].start.Unix() < e[j].start.Unix() }
