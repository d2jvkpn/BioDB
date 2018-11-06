package main

import (
	"os"
	"io"
	"log"
	"fmt"
	"strings"
	"encoding/csv"
	"net/url"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)


func main () {
	db, err := sql.Open("mysql", "hello:@/BioDB")
	if err != nil { log.Fatal (err) }
	defer db.Close()

	stmt, err := db.Prepare (`insert into Taxon 
	(taxon_id, scientific_name) values (?, ?)`)

	if err != nil { log.Fatal(err) }

	rd := csv.NewReader (os.Stdin)
	rd.Comma, rd.Comment, rd.FieldsPerRecord  = '\t', '!', -1
	rd.LazyQuotes = true

	var i int
	StartAt := time.Now()

	for {
		record, err := rd.Read ()
		if err == io.EOF { break }
		if err != nil { log.Println (err); continue }
		
		i++

		record[1] = url.QueryEscape (strings.ToLower (record[1]))
		_, err = stmt.Exec (record[0], record[1])
		if err != nil {log.Println (err)}
	}

	fmt.Println("%d records have be imported in %s", i, time.Now().Sub(StartAt))
}
