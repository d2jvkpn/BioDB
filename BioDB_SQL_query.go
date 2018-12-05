package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/d2jvkpn/gopkgs/biodb"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"os"
	"regexp"
)

const USAGE = `Query BioDB, usage:
  $ BioDB_SQL_query  <table_name>  <taxon>
    arguments:
    "Taxonomy"   <taxon_id | taxon_name>,     exactly match
    "Genome"     <taxon_id | organism name>,  exactly | ambigutily match
    "GO"         <taxon_id>,                  exactly match
    "Pathway"    <taxon_id>,                  exactly match

author: d2jvkpn
version: 0.7
release: 2018-12-05
project: https://github.com/d2jvkpn/BioDB
lisense: GPLv3 (https://www.gnu.org/licenses/gpl-3.0.en.html)
`

func main() {
	if len(os.Args) != 3 {
		fmt.Println(USAGE)
		os.Exit(2)
	}

	table, taxon := os.Args[1], os.Args[2]
	taxon = strings.Join(strings.Fields(taxon), " ")
	isdigital, _ := regexp.MatchString("^[1-9][0-9]*$", taxon)
	var err error
	var DB *sql.DB

	if DB, err = sql.Open("mysql", "hello:@/BioDB"); err != nil {
		log.Fatal(err)
	}

	defer DB.Close()

	switch {
	case table == "Taxonomy":
		var tlist []*biodb.Taxon_infor

		if tlist, err = biodb.QueryTaxonomy(DB, taxon); err != nil {
			break
		}

		jsbytes, _ := json.MarshalIndent(tlist, "", "  ")
		fmt.Println(string(jsbytes))

	case table == "Genome":
		var glist []*biodb.Genome

		if glist, err = biodb.QueryGenome(DB, taxon); err != nil {
			break
		}

		jsbytes, _ := json.MarshalIndent(glist, "", "  ")
		fmt.Println(string(jsbytes))

	case isdigital && table == "GO":
		var wt biodb.Writer
		wt = os.Stdout
		err = biodb.QueryGO(DB, taxon, wt)

	case isdigital && table == "Pathway":
		var wt biodb.Writer
		wt = os.Stdout
		err = biodb.QueryPathway(DB, taxon, wt)

	default:
		fmt.Println(USAGE)
		os.Exit(2)
	}

	if err != nil {
		log.Fatal(err)
	}
}
