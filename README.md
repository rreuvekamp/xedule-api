Xedule API
==========

This is an unofficial API for Xedule (https://summacollege.xedule.nl). Xedule is the software my school uses to schedule classes.
A friend of mine [made one](https://github.com/darkwater/xedule-api) that is more extensive, so use that. I've made this one to work just the way I like it, and of course for learning purposes. 

The only method it has is /schedule.json, which gives a week's worth of events of an attendee, ordered by the start time of the event.

##### Parameters
- aid    int (attendee id) (required)
- year   int
- week   int
- legacy type-less (other format if not empty)
- indent int (indent JSON by given amount of spaces)

For example: /schedule.json?aid=14339&week=17

##### Format
	[
		year int
		week int
		days: [
			day int // Day of week
			events: [
				start   int        // UnixTime
				end     int        // UnixTime
				desc    string     // Description
				classes [ string ]
				facs    [ string ] // Facilities
				staffs  [ string ]
			]
		]
	]

##### Features
WeekSchedules are cached in memory for 10 minutes. 

A list of attendees in the database is required to be able to put attendees at the proper type (see format; class, facility, staff).
Attendees in the database are not updated automatically. To update them, give --update-attendees with the location id, when starting the application. For example: --update-attendees=34 (for fetching and updating all attendees at location 34).
Database details can be set in the configuration file.
