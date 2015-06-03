package weeks

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rreuvekamp/xedule-api/attendee"
	"github.com/rreuvekamp/xedule-api/misc"
	"golang.org/x/net/html"
)

type week [2]int
type weeks []week

const urlWeeks = "%sAttendee/ScheduleCurrent/%d?Code=henk&attId=%d&OreId=%d" // fmt.Sprintf(urlWeeks, a.Id, a.Type, a.Lid)

// Get either fetches the weekslist from an external page (see urlWeeks) and returns is
// or returns the list from cache, if the cache is valid (and not out dated).
func Get(cache bool) (weeks, time.Time, error) {

	if cache {
		// Check for valid cache, return that if valid.
		ch := make(chan weeksResponse)
		chWksReq <- weeksRequest{ch: ch, maxAge: defReqMaxAge}
		wks := <-ch
		if wks.found {
			return wks.wks, wks.time, nil
		}
	}

	// Get the attendee for fetching the weeks list.
	chA := make(chan attendee.Attendee)
	chAttReq <- chA
	a := <-chA

	var w weeks

	// Fetch the page
	resp, err := http.Get(fmt.Sprintf(urlWeeks, misc.UrlPrefix, a.Id, a.Type, a.Lid))
	if err != nil {
		log.Println("ERROR fetching weeks:", err)
		return w, time.Time{}, err
	}
	defer resp.Body.Close()

	// Parse the HTML document
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Println("ERROR reading from fetched weeks:", err)
		return w, time.Time{}, err
	}
	w = parse(doc, w)

	sort.Sort(sort.Reverse(w))

	// Update the cache
	chWksSet <- weeksCache{wks: w, time: time.Now()}

	return w, time.Now(), nil
}

// parse is used by Get to parse and HTML node and look for
// option tags with a year/week for weeks.
// It calls itself as it goes through the nodes.
func parse(n *html.Node, w weeks) weeks {
	if n.Type == html.ElementNode && n.Data == "option" {
		var correct bool

		// Check if this option node is in a/the correct select node.
		for _, a := range n.Parent.Attr {
			if a.Key == "id" && a.Val == "currentWeek" {
				correct = true
				break
			}
		}

		// Not our select node?
		if !correct {
			goto next
		}

		var wkStr string

		// Get the value attribute's text.
		for _, a := range n.Attr {
			if a.Key == "value" {
				wkStr = a.Val
				break
			}
		}

		// Format of value attribute: year/week (e.g. 2015/19)
		var wkSplit = strings.Split(wkStr, "/")
		if len(wkSplit) < 2 {
			log.Println("WARNING: wkSplit is not 2 in weeks.parse:", wkSplit)
			goto next
		}

		yr, err1 := strconv.Atoi(wkSplit[0])
		wk, err2 := strconv.Atoi(wkSplit[1])
		if err1 != nil || err2 != nil {
			log.Println("WARNING: while atoi-ing wkSplit of fetched weeks:", err1, err2)
			goto next
		}

		w = append(w, week{yr, wk})
	}

	// Next node
next:
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		w = parse(c, w)
	}
	return w
}

// Sort functions

func (w weeks) Len() int      { return len(w) }
func (w weeks) Swap(i, j int) { w[i], w[j] = w[j], w[i] }
func (w weeks) Less(i, j int) bool {
	if len(w[i]) < 2 || len(w[j]) < 2 {
		return false
	}
	wI := strconv.Itoa(w[i][1])
	if len(wI) < 2 {
		wI = "0" + wI
	}
	wJ := strconv.Itoa(w[j][1])
	if len(wJ) < 2 {
		wJ = "0" + wJ
	}
	lessI, _ := strconv.Atoi(strconv.Itoa(w[i][0]) + wI)
	lessJ, _ := strconv.Atoi(strconv.Itoa(w[j][0]) + wJ)
	return lessI < lessJ
}
