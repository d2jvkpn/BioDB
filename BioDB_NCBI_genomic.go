package main

import (
	//"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const USAGE = `Query NCBI genomic information by providing taxonomy id, 
scientific name or genome url, usage:
  $ BioDB_NCBI_genomic  [-d directory]  <input>
`

const LISENSE = `
author: d2jvkpn
version: 0.5
release: 2018-11-23
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

	var result Genomic
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

func (result *Genomic) Save(dir string) (err error) {
	pd := path.Join(dir, strings.SplitN(result.Meta[2], "\t", 2)[1])

	err = os.MkdirAll(pd, 0755)
	if err != nil {
		return
	}

	wtj, err := os.Create(pd + "/genomic.ini")
	if err != nil {
		return
	}
	defer wtj.Close()

	fmt.Fprintln(wtj, result.ToIni())

	wts, err := os.Create(pd + "/download.sh")
	if err != nil {
		return
	}
	defer wts.Close()

	wts.Write([]byte(fmt.Sprintf(
		"#! /bin/bash\n## URL: %s\nset -eu\n",
		strings.SplitN(result.Meta[0], "\t", 2)[1]),
	))

	var ftp, cmd string
	cmd = "\n{\nwget -c -O %s -o %[1]s.dl.logging \\\n%s &&" + 
		"\nrm %[1]s.dl.logging\n} &\n"
	var fds []string

	for i, _ := range result.Reference {
		ftp = strings.SplitN(result.Reference[i], "\t", 2)[1]

		if !strings.HasPrefix(ftp, "ftp://") ||
			strings.HasSuffix(ftp, "/") {
			continue
		}

		fds = strings.SplitN(path.Base(ftp), "_", -1)

		wts.Write([]byte(fmt.Sprintf(cmd, fds[len(fds)-1], ftp)))
	}

	wts.Write([]byte("\nwait\n"))

	log.Printf("saved to %s\n", pd)
	return
}

func (result *Genomic) ToIni() (str string) {
	str = "\n[Meta]\n"
	for _, s := range result.Meta {
		str += strings.Replace(s, "\t", " = ", 1) + "\n"
	}

	str += "\n[Lineage]\n"
	for _, s := range result.Lineage {
		str += strings.Replace(s, "\t", " = ", 1) + "\n"
	}

	str += "\n[Reference]\n"
	for _, s := range result.Reference {
		str += strings.Replace(s, "\t", " = ", 1) + "\n"
	}

	str += "\n[Annotation]\n"
	for _, s := range result.Annotation {
		str += strings.Replace(s, "\t", " = ", 1) + "\n"
	}

	return
}

func (result *Genomic) Query(query string) (err error) {
	res, err := http.Get(query)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = errors.New (fmt.Sprintf("%d %s", res.StatusCode, res.Status))

		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return
	}

	if strings.Contains(doc.Find("#messagearea").Text(), "No items found") {
		err = errors.New("genome not found")

		return
	}

	result.Meta = append(result.Meta, "URL\t"+query)

	result.Meta = append(result.Meta,
		time.Now().Format("AcessTime\t2006-01-02 15:04:05 -0700"))

	result.Meta = append(result.Meta, "Name\t")

	doc.Find("span.GenomeLineage").Eq(0).Find("a").Each(func(i int,
		sel *goquery.Selection) {

		if i%2 == 0 {
			href, _ := sel.Attr("href")
			result.Lineage = append(result.Lineage,
				path.Base(href)+"\t"+sel.Text())
		}
	})

	sname := strings.Split(result.Lineage[len(result.Lineage)-1], "\t")[1]
	result.Meta[2] = result.Meta[2] + "NCBI__" + url.QueryEscape(sname)

	doc.Find("div.refgenome_sensor").Eq(0).Find("a").Each(func(i int,
		sel *goquery.Selection) {

		href, _ := sel.Attr("href")

		if strings.HasPrefix(href, "/") {
			href = "https://www.ncbi.nlm.nih.gov" + href
		}

		if sel.Text() == "genome" && strings.HasPrefix(href, "ftp:") {
			result.Meta[2] +=  "__" +
				strings.TrimRight(path.Base(href), "_genomic.fna.gz")
		}

		result.Reference = append(result.Reference, sel.Text()+"\t"+href)
	})

	count := make([]int, 6)

	doc.Find("table.GenomeList2").Find("tr").Each(func(i int,
		sel *goquery.Selection) {

		var row []string
		var s string
		var n int
		var err error

		if i == 0 {
			sel.Find("th").Each(func(i int, sel1 *goquery.Selection) {
				row = append(row, strings.TrimSpace(sel1.Text()))
			})

			result.Annotation = append(result.Annotation, row[len(row)-6:]...)
		} else {
			sel.Find("td").Each(func(i int, sel1 *goquery.Selection) {
				row = append(row, strings.TrimSpace(sel1.Text()))
			})

			for i, _ := range count {
				if len(row) < 6 {
					continue
				}

				s = strings.Replace(row[len(row)-6+i], ",", "", -1)
				if n, err = strconv.Atoi(s); err == nil {
					count[i] += n
				}
			}
		}
	})

	for i, _ := range count {
		result.Annotation[i] += "\t" + strconv.Itoa(count[i])
	}

	return
}

type Genomic struct {
	Meta, Lineage         []string
	Reference, Annotation []string
}
