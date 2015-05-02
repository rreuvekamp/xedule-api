package misc

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/ziutek/mymysql/godrv"
)

var Db *sql.DB

const UrlPrefix = "https://summacollege.xedule.nl/"

func ConnectDb() {
	var err error
	Db, err = sql.Open("mymysql", fmt.Sprintf("%s/%s/%s", "xedule-api", "xedule-api", "temppass"))

	if err != nil {
		log.Fatal(err)
		return
	}
}
