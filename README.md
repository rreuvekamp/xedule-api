Xedule API
==========

This is an unofficial API for Xedule (https://summacollege.xedule.nl). Xedule is software my school uses for scheduling classes.
A friend of mine [made one](https://github.com/darkwater/xedule-api) that is more extensive, so use that. I've made this one to work just the way I like it, and of course for learning purposes. 

It only has two methods at the moment:

### /schedule.json
Gives a week's worth of events of an attendee, ordered by the start time of the event.

##### Parameters
- aid    int (attendee id) (required)
- year   int
- week   int
- legacy type-less (if not empty, data is in a different format)
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

### /weeks.json
Gives a list of years/weeks of which there are schedules.

##### Format
	[
		[
			year int
			week int
		]

	]


#### Features
WeekSchedules are cached in memory for 10 minutes. 

Weeks list is cached in memory for 30 minutes.

A list of attendees in the database is required to be able to put attendees at the proper type (see format; class, facility, staff).
Also, weeks.json needs a valid attendee to fetch the weeks list from Xedule. The first one in the database is used for that. 

Attendees in the database are not updated automatically. To update them, give --update-attendees with the location id, when starting the application. For example: --update-attendees=34 (for fetching and updating all attendees at location 34).

Database access details can be set in the configuration file.

Database structure for table 'attendee':
	id   int primary
	name varchar(32) // Length 32 should be more than enough
	type tinyint(4)  // 1: class, 2: staff, 3: facility
	lid  int         // Location Id
