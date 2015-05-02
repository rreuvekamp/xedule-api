package wsched

import "time"

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

var legacyTimeAdd = time.Duration(time.Hour * 2)

func (w WeekSchedule) Legacy() []legacyDay {
	var days []legacyDay
	for _, d := range w.Days {
		var date string
		if len(d.Events) > 0 {
			date = d.Events[0].start.Format(legacyDate)
		}

		var events []legacyEvent
		for _, e := range d.Events {
			events = append(events, legacyEvent{
				Start:   e.start.Add(legacyTimeAdd).Format(legacyTime),
				End:     e.end.Add(legacyTimeAdd).Format(legacyTime),
				Desc:    e.Desc,
				Facs:    e.Facs,
				Staffs:  e.Staffs,
				Classes: e.Classes,
			})
		}

		days = append(days, legacyDay{
			Date:   date,
			Events: events,
		})
	}
	return days
}
