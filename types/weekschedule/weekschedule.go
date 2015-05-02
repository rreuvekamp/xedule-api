package wsched

import (
	"fmt"
	"log"
	"time"

	"github.com/PuloV/ics-golang"
)

type WeekSchedule struct {
	Days []DaySchedule
	Year int
	Week int
	//Class types.Class
}

type DaySchedule struct {
	Items []ItemSchedule
	Day   time.Weekday
}

type ScheduleItem struct {
	Start   time.Time
	End     time.Time
	Classes []string // (Other) classes/attendees
	Facs    []string // Facilities
	Staffs  []string
	Desc    string // Description
	DescH   string // Description human readable
}

const urlWeekSchedule = "https://summacollege.xedule.nl/Calendar/iCalendarICS/%d?year=%d&week=%d"

func Fetch(aid int, year int, week int) (WeekSchedule, error) {
	// ICS Parser
	p := ics.New()
	in := p.GetInputChan()

	// URL of ICS file
	in <- fmt.Sprintf(urlWeekSchedule, aid, year, week)

	// Wait for the file to be fetched and parsed.
	p.Wait()

	cal, err := p.GetCalendars()
	if err != nil || len(cal) == 0 {
		log.Println("ERROR parsing ICS files:", err, aid, year, week)
		return WeekSchedule{}, err
	}

	var days []DaySchedule

	for _, e := range cal[0].GetEvents() {

		day := e.GetStart().Weekday()

		item := ItemSchedule{
			Start: e.GetStart(),
			End:   e.GetEnd(),
			Desc:  e.GetDescription(),
			Fac:   e.GetLocation(),
		}
		fmt.Println(e.GetAttendees())

		// Check if day exists already
		exists := 0
		for i, d := range days {
			if d.Day == day {
				exists = i
				break
			}
		}
		if exists == 0 {
			days = append(days, DaySchedule{
				Day:   day,
				Items: []ItemSchedule{item},
			})
			continue
		}
		days[exists].Items = append(days[exists].Items, item)
	}

	return WeekSchedule{}, nil
}

/*
// If Schedule should be stored ('cached') in database.

type ScheduleItemAttendee struct {
	ItemScheduleId int
	AttendeeId int
}

func fetchFromDb(aid int, year int, week int) {
	// Select ScheduleItems with ScheduleItemAttendee aid is given aid.
}
*/
