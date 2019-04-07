package main

import (
	"regexp"
	"strconv"
)

type Query struct {
	Target, Term string
}

func (QF *Query) Verify() bool {
	isdigital, _ := regexp.MatchString("^\\d+$", QF.Target)

	if QF.Target != "Taxonomy" &&
		QF.Target != "Subclass" &&
		QF.Target != "Genome" &&
		QF.Target != "GO" &&
		QF.Target != "Pathway" {
		return false
	}

	if isdigital {
		n, _ := strconv.Atoi(QF.Term)
		return n > 0
	}

	return true
}
