package attendee

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/rreuvekamp/xedule-api/misc"
	"golang.org/x/net/html"
)

type Attendee struct {
	Id   int
	Name string
	Type Type
	Lid  int // locationId
}

type Type uint8

const (
	Class Type = 1 << iota // 1
	Staff Type = 1 << iota // 2
	Facil Type = 1 << iota // 3, facility
)

// FetchS: FetchSql
func FetchS(fields []string, end string) ([]Attendee, error) {

	var atts []Attendee

	// Default fieldsStr when it is empty.
	fieldsStr := "*"
	if len(fields) > 0 {
		fieldsStr = strings.Join(fields, ", ")
	}

	rows, err := misc.Db.Query("SELECT " + fieldsStr + " FROM attendee " + end + ";")
	if err != nil {
		log.Println("ERROR fetching attendee(s):", err)
		return atts, err
	}

	// Get the columns fetched.
	cols, err := rows.Columns()
	if err != nil {
		log.Println("ERROR getting columns for attendee:", err)
		return atts, err
	}
	count := len(cols)

	// Magic. Could do some cleaning.
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		var att Attendee

		// Loop over columns.
		for i, col := range cols {

			var v interface{}

			val := values[i]

			b, ok := val.([]byte)

			if ok {
				v = string(b)
			} else {
				v = val
			}

			switch col {
			case "id":
				att.Id, _ = strconv.Atoi(v.(string))
			case "name":
				att.Name = v.(string)
			case "type":
				ty, _ := strconv.Atoi(v.(string))
				att.Type = Type(ty)
			case "lid":
				att.Lid, _ = strconv.Atoi(v.(string))
			}
		}
		atts = append(atts, att)
	}
	return atts, nil
}

// Update fetches the page with attendees of the given locationId
// and saves changes in the database.
// Update is not called 'automaticly' in this program.
// It will manually be executed when it's needed.
func Update(lid int) error {

	// Fetch the page
	resp, err := http.Get(fmt.Sprintf("%s/OrganisatorischeEenheid/Attendees/%d", misc.UrlPrefix, lid))
	if err != nil {
		log.Println("ERROR fetching attendees:", err, lid)
		return err
	}
	defer resp.Body.Close()

	// Parse the document to extract the values of option tags and make a []Attendee out of that.
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Println("ERROR reading from fetched attendees:", err, lid)
		return err
	}
	atts := parse(doc, []Attendee{}, lid)

	// Save the Attendees in the database.
	for _, a := range atts {
		res, err := misc.Db.Exec(`
			INSERT INTO attendee (id, name, type, lid)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				name = VALUES(name), type = VALUES(type), lid = VALUES(lid);
		`, a.Id, a.Name, a.Type, a.Lid)
		fmt.Print(res.LastInsertId())
		fmt.Print(res.RowsAffected())
		fmt.Println(a, err)
	}

	return nil
}

// parse is used by Update to parse an HTML node.
// It calls itself while it goes through the nodes.
func parse(n *html.Node, atts []Attendee, lid int) []Attendee {
	if n.Type == html.ElementNode && n.Data == "option" {
		var rawUrl string

	attrs:
		for _, a := range n.Attr {
			if a.Key == "value" {
				rawUrl = a.Val
				break attrs
			}
		}

		// Parse the URL to get the id and name
		// Example URL:14253?Code=ZWN&attId=2&OreId=34
		u, err := url.Parse(rawUrl)
		if err != nil {
			fmt.Println("ERROR parsing URL:", err, rawUrl)
			goto next
		}

		id, err := strconv.Atoi(u.Path)
		if err != nil {
			fmt.Println("ERROR strconv of id in attendee.Update:", err)
			goto next
		}

		// TO DO: Fix variables being declared after goto statements above.

		// Get name and type
		q, _ := url.ParseQuery(u.RawQuery)

		// Attendee Type
		var ty Type
		if tys, ok := q["attId"]; ok && len(tys) > 0 {
			if tyInt, err := strconv.Atoi(tys[0]); err == nil {
				ty = Type(tyInt)
			}
		}

		// Attendee Name
		var name string
		if names, ok := q["Code"]; ok && len(names) > 0 {
			name = names[0]
		}

		// Make attendee, and append to slice
		atts = append(atts, Attendee{
			Id:   id,
			Name: name,
			Type: ty,
			Lid:  lid,
		})
	}

next:
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		atts = parse(c, atts, lid)
	}
	return atts
}

func save() {

}
