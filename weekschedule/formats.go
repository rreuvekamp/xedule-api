package wsched

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rreuvekamp/xedule-api/attendee"
)

type legacyDay struct {
	Date   string        `json:"date"`
	Events []legacyEvent `json:"events"`
}

type legacyEvent struct {
	Start   string   `json:"start"`
	End     string   `json:"end"`
	Desc    string   `json:"description"`
	Facs    []string `json:"facilities"`
	Staffs  []string `json:"staffs"`
	Classes []string `json:"classes"`
}

const legacyDate = "Mon Jan 02 2006"
const legacyTime = "15:04"

var legacyTimeAdd = time.Duration(time.Hour * 0)

// Legacy formats the WeekSchedule in a []legacyday.
func (w WeekSchedule) Legacy() []legacyDay {
	var days []legacyDay
	for _, d := range w.Days {
		var date string
		if len(d.Events) > 0 {
			date = d.Events[0].start.Format(legacyDate)
		}

		var events []legacyEvent
		for _, e := range d.Events {
			var facs, staffs, classes []string
			for _, aid := range e.Atts {
				att, ok := w.Atts[strconv.Itoa(aid)]
				if !ok {
					continue
				}

				switch att.Type {
				case attendee.Class:
					classes = append(classes, att.Name)
				case attendee.Staff:
					staffs = append(staffs, att.Name)
				case attendee.Facil:
					facs = append(facs, att.Name)
				default:
					fmt.Println("Default", att, attendee.Class, attendee.Staff, attendee.Facil)
				}
			}

			events = append(events, legacyEvent{
				Start:   e.start.Add(legacyTimeAdd).Format(legacyTime),
				End:     e.end.Add(legacyTimeAdd).Format(legacyTime),
				Desc:    e.Desc,
				Facs:    facs,
				Staffs:  staffs,
				Classes: classes,
			})
		}

		days = append(days, legacyDay{
			Date:   date,
			Events: events,
		})
	}
	return days
}
