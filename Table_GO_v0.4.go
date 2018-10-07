package main

import (
	"os"
	"io"
	"log"
	"sync"
	"strings"
	"encoding/csv"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


func main () {
	db, err := sql.Open("mysql", "hello:world@/BioDB")
	if err != nil { log.Fatal (err) }
	defer db.Close()

	stmt, err := db.Prepare (`insert into GO 
	(GO_id, prot_id, class, genes, taxon_id) values (?, ?, ?, ?, ?)`)

	if err != nil { log.Fatal(err) }

	rd := csv.NewReader (os.Stdin)
	rd.Comma, rd.Comment, rd.FieldsPerRecord  = '\t', '!', -1
	var wg sync.WaitGroup
	ch := make (chan struct{}, 50)

	for {
		record, err := rd.Read ()
		if err == io.EOF { break }
		if err != nil { log.Println (err); continue }

		rc := []string {record[4], record[7], record[8],
		record[10], strings.TrimLeft (record[12], "taxon:")}

		ch <- struct {}{}
		wg.Add (1)
		go Insert (rc, stmt, &wg, ch)
	}

	wg.Wait()
}

func Insert (rc []string, stmt *sql.Stmt, 
	wg *sync.WaitGroup, ch <- chan struct{}) {

	defer func () { <- ch; wg.Done() }()
	_, err := stmt.Exec (rc[0], rc[1], rc[2], rc[3], rc[4])
	if err != nil {log.Println (err)}
}
