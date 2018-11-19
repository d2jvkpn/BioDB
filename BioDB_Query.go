package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/d2jvkpn/gopkgs/cmdplus"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"regexp"
)

const USAGE = `Query BioDB, usage:
  $ BioDB_Query  <table_name>  <taxon>
    arguments:
    "Taxonomy"   <taxon_id | taxon_name>
    "GO"         <taxon_id>
    "Pathway"    <taxon_id>

author: d2jvkpn
version: 0.0.3
release: 2018-11-20
project: https://github.com/d2jvkpn/BioDB
lisense: GPLv3 (https://www.gnu.org/licenses/gpl-3.0.en.html)
`

type Taxon_infor struct {
	Taxon_id, Scientific_name string
	Taxon_rank, Parent_id     string
}

type Taxonomy struct {
	Taxon_infor
	Escape_name string
}

type GO struct {
	Taxon_id, Genes, GO_id string
}

type Pathway struct {
	Taxon_id, Pathway_id, Gene_id string
	KO_id, KO_information, EC_ids string
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println(USAGE)
		os.Exit(2)
	}

	table, taxon := os.Args[1], os.Args[2]
	isdigital, _ := regexp.MatchString("^[1-9][0-9]*$", taxon)
	var err error

	db, err := sql.Open("mysql", "hello:@/BioDB")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch {
	case isdigital && table == "Taxonomy":
		var t *Taxonomy
		t, err = QueryTaxonID(db, taxon)

		if err == nil {
			jsbytes, _ := json.MarshalIndent(
				[]Taxon_infor{t.Taxon_infor}, "", "  ")

			fmt.Println(string(jsbytes))
		}
	case !isdigital && table == "Taxonomy":
		err = QueryTaxonName(db, taxon)
	case isdigital && table == "GO":
		err = QueryGO(db, taxon)
	case isdigital && table == "Pathway":
		err = QueryPathway(db, taxon)
	default:
		fmt.Println(USAGE)
		os.Exit(2)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func QueryTaxonID(db *sql.DB, taxon_id string) (t *Taxonomy, err error) {
	row := db.QueryRow(
		fmt.Sprintf("select * from Taxonomy where taxon_id = '%s';", taxon_id),
	)

	t = new(Taxonomy)

	err = row.Scan(&t.Taxon_id, &t.Scientific_name, &t.Taxon_rank,
		&t.Parent_id, &t.Escape_name)

	return
}

func QueryTaxonName(db *sql.DB, taxon_name string) (err error) {
	rows, err := db.Query(
		fmt.Sprintf("select * from Taxonomy where escape_name = '%s';",
			cmdplus.NameEscape(taxon_name)),
	)

	if err != nil {
		return
	}
	defer rows.Close()

	var t Taxonomy
	var inforlist []Taxon_infor

	for rows.Next() {
		err = rows.Scan(&t.Taxon_id, &t.Scientific_name, &t.Taxon_rank,
			&t.Parent_id, &t.Escape_name)

		if err != nil {
			return
		}

		inforlist = append(inforlist, t.Taxon_infor)
	}

	if err = rows.Err(); err != nil {
		return
	}

	if len(inforlist) == 0 {
		inforlist, err = QueryTaxonHomotypic(db, taxon_name)
		if err != nil {
			return
		}
	}

	jsbytes, _ := json.MarshalIndent(inforlist, "", "  ")
	fmt.Println(string(jsbytes))

	return
}

func QueryTaxonHomotypic(db *sql.DB, taxon_name string) (
	inforlist []Taxon_infor, err error) {

	var h struct {
		id, name string
	}

	var t *Taxonomy

	rows, err := db.Query(
		fmt.Sprintf("select * from Taxonomy_homotypic where name = '%s';",
			cmdplus.NameEscape(taxon_name)),
	)

	for rows.Next() {
		err = rows.Scan(&h.id, &h.name)

		if err != nil {
			return
		}

		t, err = QueryTaxonID(db, h.id)

		if err != nil {
			return
		}

		inforlist = append(inforlist, t.Taxon_infor)
	}

	err = rows.Err()

	return
}

func QueryGO(db *sql.DB, taxon_id string) (err error) {
	rows, err := db.Query(fmt.Sprintf(
		"select * from GO where taxon_id = '%s';", taxon_id),
	)

	if err != nil {
		return
	}
	defer rows.Close()

	var t GO
	fmt.Println("genes\tGO_id")
	for rows.Next() {
		err = rows.Scan(&t.Taxon_id, &t.Genes, &t.GO_id)
		if err != nil {
			return
		}
		fmt.Printf("%s\t%s\n", t.Genes, t.GO_id)
	}

	if err = rows.Err(); err != nil {
		return
	}
	return
}

func QueryPathway(db *sql.DB, taxon_id string) (err error) {
	rows, err := db.Query(
		fmt.Sprintf("select * from Pathway where taxon_id = '%s';", taxon_id),
	)

	if err != nil {
		return
	}
	defer rows.Close()

	var t Pathway
	fmt.Println("pathway_id\tgene_id\tKO_id\tKO_information\tEC_ids")
	for rows.Next() {
		err = rows.Scan(&t.Taxon_id, &t.Pathway_id, &t.Gene_id, &t.KO_id,
			&t.KO_information, &t.EC_ids)

		if err != nil {
			return
		}

		fmt.Printf("%s\t%s\t%s\t%s\t%s\n", t.Pathway_id, t.Gene_id, t.KO_id,
			t.KO_information, t.EC_ids)
	}

	if err = rows.Err(); err != nil {
		return
	}
	return
}
