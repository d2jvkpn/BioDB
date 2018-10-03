package main

import (
	"os"
	"io"
	"log"
	"encoding/csv"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


func main () {
	db, err := sql.Open("mysql", "hello:world@/BioDB")
	if err != nil { log.Fatal (err) }
	defer db.Close()

	stmt, err := db.Prepare (`insert into Taxon 
	(taxon_id, scientific_name) values (?, ?)`)

	if err != nil { log.Fatal(err) }

	rd := csv.NewReader (os.Stdin)
	rd.Comma, rd.Comment, rd.FieldsPerRecord  = '\t', '!', -1
	rd.LazyQuotes = true

	for {
		record, err := rd.Read ()
		if err == io.EOF { break }
		if err != nil { log.Println (err); continue }

		
		_, err = stmt.Exec (record[0], record[1])
		if err != nil {log.Println (err)}
	}
}
