package main

import (
	"flag"
	"fmt"
	"github.com/d2jvkpn/gopkgs/biodb"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const USAGE = `Query NCBI genomic information by providing taxonomy id, 
scientific name or genome url, usage:
  $ BioDB_NCBI_genomic  [-d directory]  <input>
`

const LISENSE = `
author: d2jvkpn
version: 0.6
release: 2018-11-26
project: https://github.com/d2jvkpn/BioDB
lisense: GPLv3 (https://www.gnu.org/licenses/gpl-3.0.en.html)
`

func main() {
	var dir string
	flag.StringVar(&dir, "d", "", "save result to directory")

	flag.Usage = func() {
		fmt.Println(USAGE)
		flag.PrintDefaults()
		fmt.Println(LISENSE)
		os.Exit(2)
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
	}

	input := flag.Args()[0]
	prefix := "https://www.ncbi.nlm.nih.gov/genome"

	if yes, _ := regexp.MatchString("^[1-9][0-9]*$", input); yes {
		input = fmt.Sprintf(prefix+"/?term=txid%s[orgn]", input)
	} else if !strings.Contains(input, prefix) {
		input = strings.Join(strings.Fields(input), " ")
		input = prefix + "/?term=" + url.QueryEscape(input)
	}

	log.Println("querying", input)

	var result biodb.Genomic
	var err error
	if err = result.Query(input); err != nil {
		log.Fatal(err)
	}

	if dir == "" {
		fmt.Println(result.ToIni())
	} else {
		err = result.Save(dir)
	}

	if err != nil {
		log.Fatal(err)
	}
}
