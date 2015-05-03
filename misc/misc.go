package misc

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/ziutek/mymysql/godrv"
)

type config struct {
	Http struct {
		Addr string `json:"addr"`
	} `json:"http"`
	Db struct {
		User string `json:"user"`
		Pass string `json:"password"`
		Db   string `json:"database"`
	} `json:"database"`
}

const CfgFilename = "config.json"
const UrlPrefix = "https://summacollege.xedule.nl/"

var cfg config
var Db *sql.DB

// ConnectDb initializes the database 'instance'.
func ConnectDb() error {
	var err error
	Db, err = sql.Open("mymysql", fmt.Sprintf("%s/%s/%s", cfg.Db.Db, cfg.Db.User, cfg.Db.Pass)) // "xedule-api", "xedule-api", "temppass"))

	if err != nil {
		log.Println("ERROR connecting to database:", err)
		return err
	}

	return nil
}

// LoadConfig parses the configuration file with the given filename or creates one
// with empty values, if there is no file with the given filename.
func LoadConfig(filename string) error {
	data, err := ioutil.ReadFile(filename)

	// Config file with given filename doesn't exist. Create it.
	if err != nil {
		data, err := json.MarshalIndent(&cfg, "", "\t")
		if err != nil {
			log.Println("ERROR mashalling config:", err)
			return err
		}

		err = ioutil.WriteFile(filename, data, 0644)
		if err != nil {
			log.Println("ERROR writing to config file:", err, filename)
			return err
		}

		fmt.Printf("Config file created: %s. Please fill in the variables.\n", filename)
		return errors.New("config file didn't exist")
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Println("ERROR unmarshalling config file:", err, filename)
		return err
	}

	return nil
}

// Cfg returns a copy of the running configuration.
func Cfg() config {
	return cfg
}
