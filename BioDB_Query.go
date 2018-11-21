package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/d2jvkpn/gopkgs/cmdplus"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"regexp"
	"strings"
)

const USAGE = `Query BioDB, usage:
  $ BioDB_Query  <table_name>  <taxon>
    arguments:
    "Taxonomy"   <taxon_id | taxon_name>
    "GO"         <taxon_id>
    "Pathway"    <taxon_id>

author: d2jvkpn
version: 0.0.5
release: 2018-11-21
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
		var t *Taxonomy

		if t, err = QueryTaxonID(db, taxon); err == nil {
			jsbytes, _ := json.MarshalIndent(
				[]Taxon_infor{t.Taxon_infor}, "", "  ")

			fmt.Println(string(jsbytes))
		}
	case !isdigital && table == "Taxonomy":
		var inforlist []Taxon_infor

		if inforlist, err = QueryTaxonName(db, taxon); err == nil {
			jsbytes, _ := json.MarshalIndent(inforlist, "", "  ")
			fmt.Println(string(jsbytes))
		}
	case isdigital && table == "GO":
		var result [][]string

		if result, err = QueryGO(db, taxon); err == nil {
			Print2dSlice(result)
		}
	case isdigital && table == "Pathway":
		var result [][]string

		if result, err = QueryPathway(db, taxon); err == nil {
			Print2dSlice(result)
		}
	default:
		fmt.Println(USAGE)
		os.Exit(2)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func QueryTaxonID(db *sql.DB, taxon_id string) (t *Taxonomy, err error) {
	t = new(Taxonomy)

	query := fmt.Sprintf("select * from Taxonomy where taxon_id = '%s';",
		taxon_id)

	err = db.QueryRow(query).Scan(&t.Taxon_id, &t.Scientific_name,
		&t.Taxon_rank, &t.Parent_id, &t.Escape_name)

	return
}

func QueryTaxonName(db *sql.DB, taxon_name string) (inforlist []Taxon_infor,
	err error) {

	var rows *sql.Rows

	query := fmt.Sprintf("select * from Taxonomy where escape_name = '%s';",
		cmdplus.NameEscape(taxon_name, true))

	if rows, err = db.Query(query); err != nil {
		return
	}
	defer rows.Close()

	var t Taxonomy

	for rows.Next() {
		err = rows.Scan(&t.Taxon_id, &t.Scientific_name, &t.Taxon_rank,
			&t.Parent_id, &t.Escape_name)

		if err != nil {
			return
		}

		inforlist = append(inforlist, t.Taxon_infor)
	}

	if err = rows.Err(); err == nil && len(inforlist) == 0 {
		inforlist, err = QueryTaxonHomotypic(db, taxon_name)
	}

	return
}

func QueryTaxonHomotypic(db *sql.DB, taxon_name string) (
	inforlist []Taxon_infor, err error) {

	var h struct {
		id, name string
	}

	var t *Taxonomy
	var rows *sql.Rows

	query := fmt.Sprintf("select * from Taxonomy_homotypic where name = '%s';",
		cmdplus.NameEscape(taxon_name, true))

	if rows, err = db.Query(query); err != nil {
		return
	}

	for rows.Next() {
		if err = rows.Scan(&h.id, &h.name); err != nil {
			return
		}

		if t, err = QueryTaxonID(db, h.id); err != nil {
			return
		}

		inforlist = append(inforlist, t.Taxon_infor)
	}

	if err = rows.Err(); err == nil && len(inforlist) == 0 {
		err = errors.New("sql: no rows in result set")
	}

	return
}

func QueryGO(db *sql.DB, taxon_id string) (result [][]string, err error) {
	var rows *sql.Rows
	query := fmt.Sprintf("select * from GO where taxon_id = '%s';", taxon_id)

	if rows, err = db.Query(query); err != nil {
		return
	}
	defer rows.Close()

	var t GO

	result = append(result, []string{"genes", "GO_id"})

	for rows.Next() {
		if err = rows.Scan(&t.Taxon_id, &t.Genes, &t.GO_id); err != nil {
			return
		}
		result = append(result, []string{t.Genes, t.GO_id})
	}

	if err = rows.Err(); err == nil && len(result) == 1 {
		err = errors.New("sql: no rows in result set")
	}

	return
}

func QueryPathway(db *sql.DB, taxon_id string) (result [][]string, err error) {
	var rows *sql.Rows

	query := fmt.Sprintf("select * from Pathway where taxon_id = '%s';",
		taxon_id)

	if rows, err = db.Query(query); err != nil {
		return
	}
	defer rows.Close()

	var t Pathway

	result = append(result, []string{"pathway_id", "gene_id",
		"gene_information", "KO_id", "KO_information", "EC_ids"})

	for rows.Next() {
		err = rows.Scan(&t.Taxon_id, &t.Pathway_id, &t.Gene_id,
			&t.Gene_information, &t.KO_id, &t.KO_information, &t.EC_ids)

		if err != nil {
			return
		}

		result = append(result, []string{t.Pathway_id, t.Gene_id,
			t.Gene_information, t.KO_id, t.KO_information, t.EC_ids})
	}

	if err = rows.Err(); err == nil && len(result) == 1 {
		err = errors.New("sql: no rows in result set")
	}

	return
}

func Print2dSlice(result [][]string) {
	for i, _ := range result {
		fmt.Println(strings.Join(result[i], "\t"))
	}
}

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
	Taxon_id, Pathway_id, Gene_id, Gene_information string
	KO_id, KO_information, EC_ids                   string
}
