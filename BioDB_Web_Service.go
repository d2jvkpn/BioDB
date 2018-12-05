package biodb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
	"regexp"
	"strings"
)

func QueryTaxonomyID(db *sql.DB, taxon_id string) (ti *Taxon_infor, err error) {
	ti = new(Taxon_infor)
	ti.Taxon_id = taxon_id

	query := fmt.Sprintf("select scientific_name,taxon_rank,"+
		"parent_id from Taxonomy where taxon_id = '%s';", taxon_id)

	err = db.QueryRow(query).Scan(&ti.Scientific_name, &ti.Taxon_rank,
		&ti.Parent_id)

	return
}

func QueryTaxonomyName(db *sql.DB, taxon_name string) (tlist []*Taxon_infor,
	err error) {

	var rows *sql.Rows

	query := fmt.Sprintf("select taxon_id,scientific_name,taxon_rank,"+
		"parent_id from Taxonomy where escape_name = '%s';",
		NameEscape(taxon_name, true))


	if rows, err = db.Query(query); err != nil {
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

func QueryTaxonomyHomotypic(db *sql.DB, taxon_name string) (
	tlist []*Taxon_infor, err error) {

	var rows *sql.Rows

	query := fmt.Sprintf("select taxon_id from Taxonomy_homotypic where "+
		"name = '%s';", NameEscape(taxon_name, true))

	if rows, err = db.Query(query); err != nil {
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
	var query string

	isdigital, _ := regexp.MatchString("^[1-9][0-9]*$", taxon)

	if isdigital {
		query = fmt.Sprintf("select * from Genome where taxon_id = '%s';",
			taxon)
	} else {
		query = fmt.Sprintf("select * from Genome where organism_name "+
			"like '%s%%';", taxon)
	}

	if rows, err = db.Query(query); err != nil {
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

func (qf *QueryForm) MatchTaxonID(db *sql.DB) (err error) {
	var t string

	query := fmt.Sprintf("select taxon_id from %s where taxon_id = '%s';",
		qf.Table, qf.Taxon)

	err = db.QueryRow(query).Scan(&t)

	return
}

func QueryGO(db *sql.DB, taxon_id string, wt Writer) (err error) {
	var rows *sql.Rows

	query := fmt.Sprintf("select genes,GO_id from GO where taxon_id = '%s';", taxon_id)
	if rows, err = db.Query(query); err != nil {
		return
	}
	defer rows.Close()

	_, err = wt.Write([]byte("genes\tGO_id\n"))
	if err != nil {
		return
	}

	var c int
	rc := make([]string, 2)

	for rows.Next() {
		if err = rows.Scan(&rc[0], &rc[1]); err != nil {
			return
		}

		_, err = wt.Write([]byte(rc[0] + "\t" + rc[1] + "\n"))

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

func QueryPathway(db *sql.DB, taxon_id string, wt Writer) (err error) {
	var rows *sql.Rows

	query := fmt.Sprintf("select pathway_id,gene_id,gene_information,KO_id,"+
		"KO_information,EC_ids from Pathway where taxon_id = '%s';", taxon_id)

	if rows, err = db.Query(query); err != nil {
		return
	}
	defer rows.Close()

	wt.Write([]byte("pathway_id\tgene_id\tgene_information\tKO_id\t" +
		"KO_information\tEC_ids\n"))

	var c int
	rc := make([]string, 6)
	for rows.Next() {
		err = rows.Scan(&rc[0], &rc[1], &rc[2], &rc[3], &rc[4], &rc[5])
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

func NameEscape(i string, tolower bool) string {
	if tolower {
		i = strings.ToLower(i)
	}

	return url.QueryEscape(strings.Join(strings.Fields(i), " "))
}

type Writer interface {
	Write(p []byte) (n int, err error)
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
	Taxon_id, Organism_name, URL, Information string
}

type QueryForm struct {
	Taxon, Table, Download string
}
