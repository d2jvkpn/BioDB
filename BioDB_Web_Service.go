package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/d2jvkpn/gopkgs/biodb"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"regexp"
	"io/ioutil"
	"strings"
)

const USAGE = `BioDB web service, usage:
  $ BioDB_Web_Service  [-p port]
`

const LISENSE = `author: d2jvkpn
version: 0.2
release: 2018-11-30
project: https://github.com/d2jvkpn/BioDB
lisense: GPLv3 (https://www.gnu.org/licenses/gpl-3.0.en.html)
`

var db *sql.DB
var searchbts []byte

func main() {
	var err error
	var port string
	var ok bool

	flag.StringVar(&port, "p", ":8000", "set port")

	flag.Usage = func() {
		fmt.Println(USAGE)
		flag.PrintDefaults()
		fmt.Println(LISENSE)
		os.Exit(2)
	}

	flag.Parse()

	if ok, _ = regexp.MatchString("^[1-9][0-9]*$", port); ok {
		port = ":" + port
	}

	if ok, _ = regexp.MatchString("^:[1-9][0-9]*$", port); !ok {
		log.Fatalf("invalid port \"%s\"\n", port)
	}

	if db, err = sql.Open("mysql", "hello:@/BioDB"); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if searchbts, err = ioutil.ReadFile("html/search.html"); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", Search)
	http.HandleFunc("/query", Query)

	if err = http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func Search(w http.ResponseWriter, r *http.Request) {
	w.Write(searchbts)
}

func Query(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Add("StatusCode", "200")
	w.Header().Add("Status", "ok")
	w.Header().Add("Content-Type", "application/json;charset=utf-8")

	err = r.ParseForm()
	if err != nil {
		w.Header().Set("StatusCode", "400")
		w.Header().Set("Status", "invalid query")
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("invalid query"))
		return
	}

	taxon := r.FormValue("taxon")
	isdigital, _ := regexp.MatchString("^[1-9][0-9]*$", taxon)
	table := r.FormValue("table")

	switch {
	case isdigital && strings.EqualFold(table, "Taxonomy"):
		var t *biodb.Taxonomy

		if t, err = biodb.QueryTaxonID(db, taxon); err == nil {
			jsbytes, _ := json.MarshalIndent(
				[]biodb.Taxon_infor{t.Taxon_infor}, "", "  ")

			w.Write(jsbytes)
		}
	case !isdigital && strings.EqualFold(table, "Taxonomy"):
		var inforlist []biodb.Taxon_infor

		if inforlist, err = biodb.QueryTaxonName(db, taxon); err == nil {
			jsbytes, _ := json.MarshalIndent(inforlist, "", "  ")
			w.Write(jsbytes)
		}
	case isdigital && strings.EqualFold(table, "Genome"):
		var result []*biodb.Genome

		if result, err = biodb.QueryGenome(db, taxon); err == nil {

			jsbytes, _ := json.MarshalIndent(result, "", "  ")
			w.Write(jsbytes)
		}
	case isdigital && strings.EqualFold(table, "GO"):
		w.Header().Set("Content-Type", "text/plain")
		var result [][]string
		if result, err = biodb.QueryGO(db, taxon); err == nil {
			biodb.Write2dSlice(result, w)
		}
	case isdigital && strings.EqualFold(table, "Pathway"):
		w.Header().Set("Content-Type", "text/plain")
		var result [][]string
		if result, err = biodb.QueryPathway(db, taxon); err == nil {
			biodb.Write2dSlice(result, w)
		}
	default:
		w.Header().Set("StatusCode", "400")
		w.Header().Set("Status", "invalid query")
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("invalid query"))
		return
	}

	if err != nil {
		w.Header().Set("StatusCode", "404")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Status", "not found")
		w.Write([]byte(fmt.Sprintf("%s", err)))
		// w.Write([]byte("Sorry, a database error occured"))
	}
}
