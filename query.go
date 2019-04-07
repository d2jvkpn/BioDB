package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io"
	"net/url"
	"regexp"
	"strings"
)

func QueryTaxonomyID(db *sql.DB, taxon_id string) (ti *Taxon_infor, err error) {
	ti = new(Taxon_infor)
	ti.Taxon_id = taxon_id

	//query := fmt.Sprintf("select scientific_name,taxon_rank,"+
	//	"parent_id from Taxonomy where taxon_id = '%s';", taxon_id)

	err = db.QueryRow("select scientific_name,taxon_rank,parent_id "+
		"from Taxonomy where taxon_id = ?;", taxon_id).Scan(&ti.Scientific_name,
		&ti.Taxon_rank, &ti.Parent_id)

	return
}

func QueryTaxonomyName(db *sql.DB, taxon_name string) (tlist []*Taxon_infor,
	err error) {

	var rows *sql.Rows

	rows, err = db.Query("select taxon_id,scientific_name,taxon_rank,"+
		"parent_id from Taxonomy where escape_name = ?;",
		NameEscape(taxon_name, true))

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		ti := new(Taxon_infor)

		err = rows.Scan(&ti.Taxon_id, &ti.Scientific_name, &ti.Taxon_rank,
			&ti.Parent_id)

		if err != nil {
			return
		}

		tlist = append(tlist, ti)
	}

	if err = rows.Err(); err == nil && len(tlist) == 0 {
		tlist, err = QueryTaxonomyHomotypic(db, taxon_name)
	}

	return
}

func QuerySubclass(db *sql.DB, parent_id string) (
	tlist []*Taxon_infor, err error) {

	var rows *sql.Rows

	rows, err = db.Query("select taxon_id,scientific_name,taxon_rank from "+
		"Taxonomy where parent_id = ?;", parent_id)

	if err != nil {
		return
	}

	for rows.Next() {
		ti := new(Taxon_infor)
		ti.Parent_id = parent_id

		err = rows.Scan(&ti.Taxon_id, &ti.Scientific_name, &ti.Taxon_rank)

		if err != nil {
			return
		}

		tlist = append(tlist, ti)
	}

	if err = rows.Err(); err == nil && len(tlist) == 0 {
		err = sql.ErrNoRows
	}

	return
}

func QueryTaxonomyHomotypic(db *sql.DB, taxon_name string) (
	tlist []*Taxon_infor, err error) {

	var rows *sql.Rows

	rows, err = db.Query("select taxon_id from Taxonomy_homotypic where "+
		"name = ?;", NameEscape(taxon_name, true))

	if err != nil {
		return
	}

	for rows.Next() {
		ti := new(Taxon_infor)

		if err = rows.Scan(&ti.Taxon_id); err != nil {
			return
		}

		if ti, err = QueryTaxonomyID(db, ti.Taxon_id); err != nil {
			return
		}

		tlist = append(tlist, ti)
	}

	if err = rows.Err(); err == nil && len(tlist) == 0 {
		err = sql.ErrNoRows
	}

	return
}

func QueryTaxonomy(db *sql.DB, taxon string) (
	tlist []*Taxon_infor, err error) {

	isdigital, _ := regexp.MatchString("^[1-9][0-9]*$", taxon)

	if isdigital {
		var t *Taxon_infor

		if t, err = QueryTaxonomyID(db, taxon); err == nil {
			tlist = append(tlist, t)
		}

		return
	}

	tlist, err = QueryTaxonomyName(db, taxon)

	return
}

func QueryGenome(db *sql.DB, taxon string) (result []*Genome, err error) {
	var rows *sql.Rows

	rows, err = db.Query("select * from Genome where taxon_id = ?;", taxon)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		g := new(Genome)

		if err = rows.Scan(&g.Taxon_id, &g.Organism_name, &g.URL,
			&g.Information); err != nil {
			return
		}

		result = append(result, g)
	}

	if err = rows.Err(); err == nil && len(result) == 0 {
		err = sql.ErrNoRows
	}

	return
}

func (qf *Query) MatchTaxonID(db *sql.DB) (err error) {
	var t string

	query := fmt.Sprintf("select taxon_id from %s where taxon_id = ?;",
		qf.Target)

	err = db.QueryRow(query, qf.Term).Scan(&t)

	return
}

func QueryGO(db *sql.DB, taxon_id string, wt io.Writer) (err error) {
	var rows *sql.Rows

	//rows, err = db.Query("select genes,GO_id from GO where taxon_id = ?;",
	//	taxon_id)

	rows, err = db.Query("select GO.genes,GO.GO_id,GO_definition.name,"+
		"GO_definition.namespace from GO inner join GO_definition on "+
		"GO.GO_id = GO_definition.GO_id where GO.taxon_id = ?;", taxon_id)

	if err != nil {
		return
	}
	defer rows.Close()

	//_, err = wt.Write([]byte("genes\tGO_id\n"))
	_, err = wt.Write([]byte("genes\tGO_id\tname\tnamespace\n"))

	if err != nil {
		return
	}

	var c int
	// rc := make([]string, 2)
	rc := make([]string, 4)

	for rows.Next() {
		err = rows.Scan(&rc[0], &rc[1], &rc[2], &rc[3])
		//err = rows.Scan(&rc[0], &rc[1])
		if err != nil {
			return
		}

		_, err = wt.Write([]byte(strings.Join(rc, "\t") + "\n"))

		if err != nil {
			return
		}

		c++
	}

	if err = rows.Err(); err == nil && c == 0 {
		err = sql.ErrNoRows
	}

	return
}

func QueryPathway(db *sql.DB, taxon_id string, wt io.Writer) (err error) {
	var rows *sql.Rows

	//rows, err = db.Query("select pathway_id,gene_id,gene_information,KO_id,"+
	//	"KO_information,EC_ids from Pathway where taxon_id = ?;", taxon_id)

	rows, err = db.Query("select Pathway.Pathway_id,Pathway.gene_id,"+
		"Pathway.gene_information,Pathway.KO_id,Pathway.KO_information,"+
		"Pathway.EC_ids,"+
		"Pathway_definition.C_id,Pathway_definition.C_name,"+
		"Pathway_definition.B_id,Pathway_definition.B_name,"+
		"Pathway_definition.A_id,Pathway_definition.A_name "+
		"from Pathway inner join Pathway_definition on "+
		"concat('C', right(Pathway.Pathway_id, 5)) = Pathway_definition.C_id "+
		"where Pathway.taxon_id = ?;", taxon_id)

	if err != nil {
		return
	}
	defer rows.Close()

	//wt.Write([]byte("pathway_id\tgene_id\tgene_information\tKO_id\t" +
	//	"KO_information\tEC_ids\n"))

	wt.Write([]byte("pathway_id\tgene_id\tgene_information\tKO_id\t" +
		"KO_information\tEC_ids\tC_id\tC_name\tB_id\tB_name\tA_id\tA_name\n"))

	var c int
	rc := make([]string, 12)

	for rows.Next() {
		err = rows.Scan(&rc[0], &rc[1], &rc[2], &rc[3], &rc[4], &rc[5], &rc[6],
			&rc[7], &rc[8], &rc[9], &rc[10], &rc[11])
		if err != nil {
			return
		}

		_, err = wt.Write([]byte(strings.Join(rc, "\t") + "\n"))
		if err != nil {
			return
		}

		c++
	}

	if err = rows.Err(); err == nil && c == 0 {
		err = sql.ErrNoRows
	}

	return
}

func NameEscape(i string, tolower bool) (s string) {
	if tolower {
		i = strings.ToLower(i)
	}

	s = strings.Replace(strings.Join(strings.Fields(i), " "), "\"\"", "\"", -1)

	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		k := []byte(s)
		s = string(k[1 : len(k)-1])
	}

	if strings.Contains(s, " ") {
		s = url.QueryEscape(s)
	}

	return
}

type Taxon_infor struct {
	Taxon_id, Scientific_name string
	Taxon_rank, Parent_id     string
}

type Taxonomy struct {
	Taxon_infor
	Escape_name string
}

type Genome struct {
	Taxon_id, Organism_name, Information string
	URL                                  template.URL
}
