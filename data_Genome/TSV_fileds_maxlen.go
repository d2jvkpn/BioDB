package main

import (
	"fmt"
	"github.com/d2jvkpn/gopkgs/cmdplus"
	"log"
	"os"
	"strconv"
	"strings"
)

const USAGE = `Calculate fileds' max length of a tsv file, usage:
  $ TSV_fileds_maxlen  <tsv_file>
`

const LISENSE = `author: d2jvkpn
version: 0.1
release: 2018-11-30
project: https://github.com/d2jvkpn/DataAnalysis
lisense: GPLv3 (https://www.gnu.org/licenses/gpl-3.0.en.html)
`

func main() {
	if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println(USAGE)
		fmt.Println(LISENSE)
		os.Exit(2)
	}

	scanner, file, err := cmdplus.ReadCmdInput(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var i int
	var hfds []string
	var fds []string
	ml := make(map[int]int)

	for scanner.Scan() {
		i++
		if i == 1 {
			hfds = strings.Split(scanner.Text(), "\t")
			for j, _ := range hfds {
				ml[j] = 0
			}
		} else {
			fds = strings.Split(scanner.Text(), "\t")

			if len(fds) != len(ml) {
				log.Fatalf("line with different number of fields at %d\n", i)
			}
			for j, _ := range fds {
				if len(fds[j]) > ml[j] {
					ml[j] = len(fds[j])
				}
			}
		}
	}

	var rc [][]string
	rc = append(rc, []string{"FIELD", "NAME", "MAX_LENGTH"})

	for j, _ := range hfds {
		rc = append(rc, []string{strconv.Itoa(j + 1), hfds[j],
			strconv.Itoa(ml[j])})
	}

	cmdplus.PrintStringArray(rc)
}
