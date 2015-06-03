package lastupdate

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rreuvekamp/xedule-api/attendee"
	"github.com/rreuvekamp/xedule-api/misc"
	"golang.org/x/net/html"
)

const urlLastUpdate = "%sAttendee/ChangeWeek/%d?Code=henk&attId=%d&OreId=%d"
const tmLayout = "2-1-2006 15:04:05"

func Get(year, week int, cache bool) (time.Time, time.Time, error) {

	if cache {
		// Check for cache
		ch := make(chan cacheResp)

		chLuReq <- cacheReq{
			ch:   ch,
			year: year,
			week: week,
		}

		c := <-ch

		// Serve cache if it exists.
		if c.found {
			fmt.Println("Found cache")
			return c.tmLu, c.tm, nil
		}
	}

	// Get the attendee for fetching the weeks list.
	chA := make(chan attendee.Attendee)
	chAttReq <- chA
	a := <-chA

	// Fetch page.
	resp, err := http.PostForm(fmt.Sprintf(urlLastUpdate, misc.UrlPrefix, a.Id, a.Type, a.Lid),
		url.Values{"currentWeek": {strconv.Itoa(year) + "/" + strconv.Itoa(week)}})
	if err != nil {
		log.Println("ERROR fetching page with last update:", err, year, week, a)
		return time.Time{}, time.Time{}, err
	}
	defer resp.Body.Close()

	// Parse page
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Println("ERROR parsing fetched last update:", err, year, week, a)
	}
	tmLu, _ := parse(doc)

	// Save cache
	chLuAdd <- cacheAdd{
		tmLu: tmLu,
		year: year,
		week: week,
		tm:   time.Now(),
	}

	return tmLu, time.Time{}, err
}

func parse(n *html.Node) (time.Time, bool) {
	if n.Type == html.ElementNode && n.Data == "div" {

		var correct bool

		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "dateCreated" {
				correct = true
				break
			}
		}

		if !correct {
			goto next
		}

		c := n.FirstChild
		if c == nil {
			goto next
		}

		split := strings.Split(c.Data, ":\n")
		if len(split) < 2 {
			log.Println("WARNING/ERROR: Len of splitted data of lastupdate is < 2:", split)
			return time.Time{}, true
		}

		str := strings.TrimSpace(split[1])

		tm, err := time.ParseInLocation(tmLayout, str, misc.Loc)
		if err != nil {
			log.Println("ERROR when parsing time of fetched lastupdate:", err)
			return time.Time{}, true
		}

		fmt.Println(tm.Unix())

		return tm, true
	}

next:
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		tm, done := parse(c)
		if done {
			return tm, done
		}
	}

	return time.Time{}, false
}
