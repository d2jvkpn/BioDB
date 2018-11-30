package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/d2jvkpn/gopkgs/biodb"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const USAGE = `BioDB web service, usage:
  $ BioDB_Web_Service  [-p port]
`

const LISENSE = `author: d2jvkpn
version: 0.3
release: 2018-11-30
project: https://github.com/d2jvkpn/BioDB
lisense: GPLv3 (https://www.gnu.org/licenses/gpl-3.0.en.html)
`

var db *sql.DB
var searchbts []byte
var invalidtmpl *template.Template
var QF biodb.QueryForm
var err error

func main() {
	var port string
	var ok bool

	defer db.Close()

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

	http.HandleFunc("/", WriteSearch)
	http.HandleFunc("/query", QueryTable)

	if err = http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func init() {
	if db, err = sql.Open("mysql", "hello:@/BioDB"); err != nil {
		log.Fatal(err)
	}

	if searchbts, err = ioutil.ReadFile("HTML/search.html"); err != nil {
		log.Fatal(err)
	}

	if invalidtmpl, err = template.ParseFiles("HTML/invalid.tmpl"); err != nil {
		log.Fatal(err)
	}
}

func WriteSearch(w http.ResponseWriter, r *http.Request) {
	w.Write(searchbts)
}

func QueryTable(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Add("StatusCode", "200")
	w.Header().Add("Status", "ok")
	w.Header().Add("Content-Type", "application/json;charset=utf-8")

	err = r.ParseForm()
	if err != nil {
		w.Header().Set("StatusCode", "400")
		w.Header().Set("Status", "invalid query")
		w.Header().Set("Status", "not found")
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		invalidtmpl.Execute(w, QF)
		return
	}

	QF.Taxon = r.FormValue("taxon")
	isdigital, _ := regexp.MatchString("^[1-9][0-9]*$", QF.Taxon)
	QF.Table = r.FormValue("table")

	switch {
	case isdigital && strings.EqualFold(QF.Table, "Taxonomy"):
		var t *biodb.Taxonomy

		if t, err = biodb.QueryTaxonID(db, QF.Taxon); err == nil {
			jsbytes, _ := json.MarshalIndent(
				[]biodb.Taxon_infor{t.Taxon_infor}, "", "  ")

			w.Write(jsbytes)
		}
	case !isdigital && strings.EqualFold(QF.Table, "Taxonomy"):
		var inforlist []biodb.Taxon_infor

		if inforlist, err = biodb.QueryTaxonName(db, QF.Taxon); err == nil {
			jsbytes, _ := json.MarshalIndent(inforlist, "", "  ")
			w.Write(jsbytes)
		}
	case isdigital && strings.EqualFold(QF.Table, "Genome"):
		var result []*biodb.Genome

		if result, err = biodb.QueryGenome(db, QF.Taxon); err == nil {
			jsbytes, _ := json.MarshalIndent(result, "", "  ")
			w.Write(jsbytes)
		}
	case isdigital && strings.EqualFold(QF.Table, "GO"):
		var bts []byte
		dispo := "inline; filename=\"Gene_Ontology.%s.tsv.gz\"" 

		if bts, err = biodb.QueryGO(db, QF.Taxon); err == nil {
			w.Header().Set("Content-Type", "application/x-gzip")
			w.Header().Set("Content-Disposition", fmt.Sprintf(dispo, QF.Taxon))

			w.Write(bts)
		}
	case isdigital && strings.EqualFold(QF.Table, "Pathway"):
		var bts []byte
		dispo := "inline; filename=\"KEGG_Pathway.%s.tsv.gz\"" 

		if bts, err = biodb.QueryPathway(db, QF.Taxon); err == nil {
			w.Header().Set("Content-Type", "application/x-gzip")
			w.Header().Set("Content-Disposition", fmt.Sprintf(dispo, QF.Taxon))

			w.Write(bts)
		}
	default:
		w.Header().Set("StatusCode", "400")
		w.Header().Set("Status", "invalid query")
		w.Header().Set("Status", "not found")
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		invalidtmpl.Execute(w, QF)
		return
	}

	if err != nil {
		w.Header().Set("StatusCode", "404")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Status", "not found")
		w.Write([]byte(fmt.Sprintf("%s", err)))
	}
}
