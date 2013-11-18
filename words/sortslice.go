package main

import (
	"sort"
)

type By func(s1, s2 string) bool

type stringSorter struct {
	words []string
	by    By
}

func (by By) Sort(s []string) {
	sS := &stringSorter{
		words: s,
		by:    by,
	}
	sort.Sort(sS)
}

func (s *stringSorter) Len() int {
	return len(s.words)
}

func (s *stringSorter) Swap(i, j int) {
	s.words[i], s.words[j] = s.words[j], s.words[i]
}

func (s *stringSorter) Less(i, j int) bool {
	return s.by(s.words[i], s.words[j])
}

func sortList(s []string) {
	length := func(s1, s2 string) bool {
		return len(s1) < len(s2)
	}

	By(length).Sort(s)
}
