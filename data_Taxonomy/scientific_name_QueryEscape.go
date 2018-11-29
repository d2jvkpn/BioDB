package main

import (
	"bufio"
	"log"
	"net/url"
	"os"
	"strings"
)

func main() {
	frd, err := os.Open("Taxonomy.0.tsv")
	if err != nil {
		log.Fatal(err)
	}
	defer frd.Close()
	r := bufio.NewScanner(frd)

	if _, err := os.Stat("Taxonomy.tsv"); !os.IsNotExist(err) {
		err := os.Remove("Taxonomy.tsv")
		if err != nil {
			log.Fatal(err)
		}
	}

	fwt, err := os.Create("Taxonomy.tsv")
	if err != nil {
		log.Fatal(err)
	}
	defer fwt.Close()

	var l, nc string
	r.Scan()
	l = r.Text()
	fwt.Write([]byte(l + "\tescape_name\n"))

	for r.Scan() {
		l = r.Text()
		nc = strings.ToLower(strings.Split(l, "\t")[1])
		fwt.Write([]byte(l + "\t" + url.QueryEscape(nc) + "\n"))
	}
}
