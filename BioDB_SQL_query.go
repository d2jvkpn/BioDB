package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/d2jvkpn/gopkgs/biodb"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"regexp"
)

const USAGE = `Query BioDB, usage:
  $ BioDB_SQL_query  <table_name>  <taxon>
    arguments:
    "Taxonomy"   <taxon_id | taxon_name>
    "GO"         <taxon_id>
    "Pathway"    <taxon_id>

author: d2jvkpn
version: 0.6
release: 2018-11-26
project: https://github.com/d2jvkpn/BioDB
lisense: GPLv3 (https://www.gnu.org/licenses/gpl-3.0.en.html)
`

func main() {
	if len(os.Args) != 3 {
		fmt.Println(USAGE)
		os.Exit(2)
	}

	table, taxon := os.Args[1], os.Args[2]
	isdigital, _ := regexp.MatchString("^[1-9][0-9]*$", taxon)
	var err error
	var db *sql.DB

	if db, err = sql.Open("mysql", "hello:@/BioDB"); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	switch {
	case isdigital && table == "Taxonomy":
		var t *biodb.Taxonomy

		if t, err = biodb.QueryTaxonID(db, taxon); err == nil {
			jsbytes, _ := json.MarshalIndent(
				[]biodb.Taxon_infor{t.Taxon_infor}, "", "  ")

			fmt.Println(string(jsbytes))
		}
	case !isdigital && table == "Taxonomy":
		var inforlist []biodb.Taxon_infor

		if inforlist, err = biodb.QueryTaxonName(db, taxon); err == nil {
			jsbytes, _ := json.MarshalIndent(inforlist, "", "  ")
			fmt.Println(string(jsbytes))
		}
	case isdigital && table == "GO":
		var result [][]string

		if result, err = biodb.QueryGO(db, taxon); err == nil {
			biodb.Write2dSlice(result, os.Stdout)
		}
	case isdigital && table == "Pathway":
		var result [][]string

		if result, err = biodb.QueryPathway(db, taxon); err == nil {
			biodb.Write2dSlice(result, os.Stdout)
		}
	default:
		fmt.Println(USAGE)
		os.Exit(2)
	}

	if err != nil {
		log.Fatal(err)
	}
}
