Xedule API
==========

This is an unofficial API for Xedule (https://summacollege.xedule.nl). Xedule is software my school uses for scheduling classes.
A friend of mine [made one](https://github.com/darkwater/xedule-api) that is more extensive, so use that. I've made this one to work just the way I like it, and of course for learning purposes. 

It has three methods at the moment:

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
	{
		year int
		week int
		days {
			day int // Day of week
			baseuts int // UnixTime of start of this day (Note: 0:00 AM, not start of first event)
			events {
				start int    // Seconds from start of day (baseuts) till start of event.
				end   int    // Seconds from start of day (baseuts) till end of event.
				desc  string // Description
				atts  int	 // Attendee ids
			}
		}
		atts { // Attendees which attent (and ids are in) one or more events. 
			id   int
			name string
			type int // 1: Class, 2: Staff, 3: Facility
		}
	}

### /weeks.json
Gives a list of years/weeks of which there are schedules.

#### Paramater
- nocache\*

##### Format
	[
		[
			int // Year
			int // Week
		]
	]

### /attendee.json
Gives attendee information of which ids are given, or attendees of the given location id.

##### Parameters
- aid int (or "int,int...") (required or) // Attendee id
- lid int (required or) // Location id
- nocache\*

For example: /attendee.json?aid=14339&aid=13451,13452

##### Format
	[
		{
			id   int
			name string
			type int // 1: Class, 2: Staff, 3: Facility
		} 
	]

### /lastupdate.json
Gives the time the schedule with the given year/week was last updated.

##### Parameters
- year int
- week int
- nocache\*

##### Format
	{
		year int
		week int
		uts int // UnixTimeStamp
	}


### Features and extra information

##### Database
A list of attendees in the database is required for /attendee.json .
The application will work without it, although other API methods also depend on it and will behave odd.

Attendees in the database are not updated automatically. To update them, give --update-attendees with the location id on the application. 
For example: --update-attendees=34 (for fetching and updating all attendees at location 34).

Database access details can be set in the configuration file which is generated on first startup.

SQL database structure for table 'attendee':
- id   int primary
- name varchar(32) // Length 32 should be more than enough
- type tinyint(4)  // 1: class, 2: staff, 3: facility
- lid  int         // Location Id

##### Cache
WeekSchedules are cached in memory for 10 minutes. 
Weeks list is cached in memory for 30 minutes.
LastUpdate times are cached in memory for 10 minutes.
* When passing nocache GET parameter, and the remote IP address is white listed for it, cache will not be looked up, thus guaranteeing up to date information. IP addresses on which nocache should have this effect can be set ('white listed') in the configuration file.
