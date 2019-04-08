package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func Search(c *gin.Context, DB *sql.DB) {
	// POST parse
	query := Query{c.PostForm("target"), c.PostForm("term")}

	// GET parse
	if query.Target == "" {
		query.Target = c.Query("target")
	}

	if query.Term == "" {
		query.Term = c.Query("term")
	}

	var err error
	var data interface{}
	var tmpl string

	if ok := query.Verify(); !ok {
		c.HTML(http.StatusBadRequest, "InvalidQuery.html", query)
		return
	}

	switch {
	case query.Target == "Taxonomy":
		data, err = GetTaxonomy(query, DB)
		tmpl = "Taxonomy.html"

	case query.Target == "Subclass":
		data, err = GetSubclass(query, DB)
		tmpl = "Taxonomy.html"

	case query.Target == "Genome":
		data, err = GetGenome(query, DB)
		tmpl = "Genome.html"

	case query.Target == "GO" || query.Target == "Pathway":
		data, err = GetGO_Pathway(query, DB)
		tmpl = "Download.html"

	default:
		c.HTML(http.StatusBadRequest, "InvalidQuery.html", &query)
		return
	}

	switch err {
	case nil:
		c.HTML(http.StatusOK, tmpl, &data)

	case sql.ErrNoRows:
		c.HTML(http.StatusNotFound, "NotFound.html", &query)

	default:
		s := "an error ocurred quering %s in %s: %s\n"
		log.Printf(s, query.Target, query.Term, err)
		c.HTML(http.StatusInternalServerError, "InternalError.html", &query)
	}

	return
}

func API(c *gin.Context, DB *sql.DB) {
	query := Query{c.Query("target"), c.Query("term")}
	var err error
	var data interface{}

	if ok := query.Verify(); !ok {
		fmt.Println(query)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	switch {
	case query.Target == "Taxonomy":
		data, err = GetTaxonomy(query, DB)

	case query.Target == "Subclass":
		data, err = GetSubclass(query, DB)

	case query.Target == "Genome":
		data, err = GetGenome(query, DB)

	case query.Target == "GO" || query.Target == "Pathway":
		data, err = GetGO_Pathway(query, DB)

	default:
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	switch err {
	case nil:
		c.JSON(http.StatusOK, &data)

	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound, nil)

	default:
		s := "an error ocurred quering %s in %s: %s\n"
		log.Printf(s, query.Target, query.Term, err)
		c.JSON(http.StatusInternalServerError, nil)
	}

	return
}

func Download(c *gin.Context, DB *sql.DB) {
	query := Query{c.Query("target"), c.Query("term")}
	var err error
	ok := query.Verify()

	if !ok || (query.Target != "GO" && query.Target != "Pathway") {
		c.HTML(http.StatusBadRequest, "InvalidQuery.html", query)
		return
	}

	var buf bytes.Buffer
	var wt io.Writer
	var dispo string
	gzw := gzip.NewWriter(&buf)
	wt = gzw

	if query.Target == "GO" {
		err = QueryGO(DB, query.Term, wt)
		dispo = "attachment; filename=\"Gene_Ontology.%s.tsv.gz\""
	} else {
		err = QueryPathway(DB, query.Term, wt)
		dispo = "attachment; filename=\"KEGG_Pathway.%s.tsv.gz\""
	}

	gzw.Close()

	switch err {
	case nil:
		c.Header("Content-Disposition", fmt.Sprintf(dispo, query.Term))
		c.Data(http.StatusOK, "application/x-gzip", buf.Bytes())
		return

	case sql.ErrNoRows:
		c.String(http.StatusNotFound, "Not Found")

	default:
		s := "an error ocurred quering %s in %s: %s\n"
		log.Printf(s, query.Target, query.Term, err)

		c.String(http.StatusInternalServerError, "service error")
	}

	return
}

func GetTaxonomy(query Query, DB *sql.DB) (data interface{}, err error) {
	var tlist []*Taxon_infor

	if tlist, err = QueryTaxonomy(DB, query.Term); err != nil {
		return
	}

	data = struct {
		*Query
		Taxonlist []*Taxon_infor
	}{&query, tlist}

	return
}

func GetSubclass(query Query, DB *sql.DB) (data interface{}, err error) {
	var tlist []*Taxon_infor

	if tlist, err = QuerySubclass(DB, query.Term); err != nil {
		return
	}

	data = struct {
		*Query
		Taxonlist []*Taxon_infor
	}{&query, tlist}

	return
}

func GetGenome(query Query, DB *sql.DB) (data interface{}, err error) {
	var glist []*Genome

	if glist, err = QueryGenome(DB, query.Term); err != nil {
		return
	}

	data = struct {
		*Query
		Genomelist []*Genome
	}{&query, glist}

	return
}

func GetGO_Pathway(query Query, DB *sql.DB) (data interface{}, err error) {
	var ti *Taxon_infor

	if err = query.MatchTaxonID(DB); err != nil {
		return
	}

	if ti, err = QueryTaxonomyID(DB, query.Term); err != nil {
		return
	}
	data = struct {
		*Query
		Scientific_name *string
	}{&query, &ti.Scientific_name}

	return
}
