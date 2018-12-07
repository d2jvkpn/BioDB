package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/d2jvkpn/gopkgs/biodb"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const USAGE = `BioDB web service, usage:
  $ BioDB_Web_Service  [-p port]
`

const LISENSE = `
author: d2jvkpn
version: 0.6
release: 2018-12-08
project: https://github.com/d2jvkpn/BioDB
lisense: GPLv3 (https://www.gnu.org/licenses/gpl-3.0.en.html)
`

var (
	DB    *sql.DB
	Tmpls map[string]*template.Template
)

const (
	DBuser   = "world"
	DBpasswd = ""
	DBhost   = "tcp(localhost:3306)"
)

func main() {
	var port string
	var err error

	flag.StringVar(&port, "p", ":8000", "set port")

	flag.Usage = func() {
		fmt.Println(USAGE)
		flag.PrintDefaults()
		fmt.Println(LISENSE)
		os.Exit(2)
	}

	flag.Parse()
	ValidPort(&port)

	DB, err = sql.Open("mysql",
		fmt.Sprintf("%s:%s@%s/BioDB", DBuser, DBpasswd, DBhost))

	if err != nil {
		log.Fatal(err)
	}

	defer DB.Close()

	http.HandleFunc("/", ServerSearch)
	http.HandleFunc("/query", QueryTable)

	if err = http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func init() {
	var err error

	Tmpls = make(map[string]*template.Template)

	Tmpls["search"], err = template.ParseFiles("html/Search.html")
	if err != nil {
		log.Fatal(err)
	}

	Tmpls["invalid"], err = template.ParseFiles("html/InvalidQuery.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	Tmpls["notfound"], err = template.ParseFiles("html/NotFound.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	Tmpls["error"], err = template.ParseFiles("html/InternalError.html")
	if err != nil {
		log.Fatal(err)
	}

	Tmpls["download"], err = template.ParseFiles("html/Download.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	var b []byte

	if b, err = ioutil.ReadFile("html/Genome.tmpl"); err != nil {
		log.Fatal(err)
	}

	Tmpls["genome"] = template.Must(template.New("").Funcs(funcMap).
		Parse(string(b)))
	// Tmpls["genome"], err = template.ParseFiles("html/Results_Genome.tmpl")

	if b, err = ioutil.ReadFile("html/Taxonomy.tmpl"); err != nil {
		log.Fatal(err)
	}

	Tmpls["taxon"] = template.Must(template.New("").Funcs(funcMap).
		Parse(string(b)))
}

func ValidPort(port *string) {
	var ok bool

	if ok, _ = regexp.MatchString("^[1-9][0-9]*$", *port); ok {
		*port = ":" + *port
		return
	}

	if ok, _ = regexp.MatchString("^:[1-9][0-9]*$", *port); !ok {
		log.Fatalf("invalid port \"%s\"\n", *port)
	}

	return
}

func ServerSearch(w http.ResponseWriter, r *http.Request) {
	Tmpls["search"].Execute(w, nil)
}

func QueryTable(w http.ResponseWriter, r *http.Request) {
	var err error
	var isdigital bool
	var ok bool
	r.ParseForm()

	QF := biodb.QueryForm{
		strings.Join(strings.Fields(r.FormValue("taxon")), " "),
		r.FormValue("table"),
		r.FormValue("download")}

	if QF.Download != "true" {
		// fmt.Println(QF.Download)
		QF.Download = "false"
	}

	if isdigital, ok = QF.IsValid(); !ok {
		w.Header().Add("StatusCode", strconv.Itoa(http.StatusBadRequest))
		w.Header().Add("Status", http.StatusText(400))
		Tmpls["invalid"].Execute(w, &QF)
		return
	}

	switch {
	case QF.Table == "Taxonomy":
		var tlist []*biodb.Taxon_infor

		if tlist, err = biodb.QueryTaxonomy(DB, QF.Taxon); err != nil {
			break
		}

		if QF.Download == "true" {
			w.Header().Add("StatusCode", "200")
			w.Header().Add("Status", "ok")
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			jsbytes, _ := json.MarshalIndent(tlist, "", "  ")
			_, err = w.Write(jsbytes)

		} else {
			data := struct {
				*biodb.QueryForm
				Taxonlist []*biodb.Taxon_infor
			}{&QF, tlist}

			err = Tmpls["taxon"].Execute(w, &data)
		}

	case isdigital && QF.Table == "Subclass":
		var tlist []*biodb.Taxon_infor

		if tlist, err = biodb.QuerySubclass(DB, QF.Taxon); err != nil {
			break
		}

		if QF.Download == "true" {
			w.Header().Add("StatusCode", "200")
			w.Header().Add("Status", "ok")
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			jsbytes, _ := json.MarshalIndent(tlist, "", "  ")
			_, err = w.Write(jsbytes)

		} else {
			data := struct {
				*biodb.QueryForm
				Taxonlist []*biodb.Taxon_infor
			}{&QF, tlist}

			err = Tmpls["taxon"].Execute(w, &data)
		}

	case isdigital && QF.Table == "Genome":
		var glist []*biodb.Genome

		if glist, err = biodb.QueryGenome(DB, QF.Taxon); err != nil {
			break
		}

		if QF.Download == "true" {
			w.Header().Add("StatusCode", "200")
			w.Header().Add("Status", "ok")
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			jsbytes, _ := json.MarshalIndent(glist, "", "  ")
			_, err = w.Write(jsbytes)

		} else {
			data := struct {
				*biodb.QueryForm
				Genomelist []*biodb.Genome
			}{&QF, glist}

			err = Tmpls["genome"].Execute(w, &data)
		}

	case isdigital && (QF.Table == "GO" || QF.Table == "Pathway"):
		if QF.Download == "true" {
			var buf bytes.Buffer
			var wt io.Writer
			var dispo string
			gzw := gzip.NewWriter(&buf)
			wt = gzw

			if QF.Table == "GO" {
				err = biodb.QueryGO(DB, QF.Taxon, wt)
				dispo = "attachment; filename=\"Gene_Ontology.%s.tsv.gz\""
			} else {
				err = biodb.QueryPathway(DB, QF.Taxon, wt)
				dispo = "attachment; filename=\"KEGG_Pathway.%s.tsv.gz\""
			}

			gzw.Close()

			if err != nil {
				break
			}

			w.Header().Set("Content-Type", "application/x-gzip")
			w.Header().Set("Content-Disposition", fmt.Sprintf(dispo, QF.Taxon))

			_, err = w.Write(buf.Bytes())

		} else {
			if err = QF.MatchTaxonID(DB); err != nil {
				break
			}

			var ti *biodb.Taxon_infor
			ti, _ = biodb.QueryTaxonomyID(DB, QF.Taxon)

			data := struct {
				Scientific_name *string
				*biodb.QueryForm
			}{&ti.Scientific_name, &QF}

			err = Tmpls["download"].Execute(w, &data)
		}

	default:
		w.Header().Add("StatusCode", strconv.Itoa(http.StatusBadRequest))
		w.Header().Add("Status", http.StatusText(400))
		Tmpls["invalid"].Execute(w, &QF)
		return
	}

	switch err {
	case nil:
		return

	case sql.ErrNoRows:
		w.Header().Add("StatusCode", strconv.Itoa(http.StatusNotFound))
		w.Header().Add("Status", http.StatusText(404))

		Tmpls["notfound"].Execute(w, &QF)

	default:
		log.Printf("an error ocurred quering %s in %s: %s\n",
			QF.Taxon, QF.Table, err)

		w.Header().Add("StatusCode",
			strconv.Itoa(http.StatusInternalServerError))

		w.Header().Add("Status", http.StatusText(500))

		Tmpls["error"].Execute(w, &QF)
	}
}

func Add(a, b int) string {
	return strconv.Itoa(a + b)
}

var funcMap = template.FuncMap{
	"Add": Add,
}
