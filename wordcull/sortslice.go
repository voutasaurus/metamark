package wordcull

import (
	"sort"
)

// By allows us to choose what functions to sort with.
type By func(s1, s2 string) bool

// stringSorter allows us to use a particular function to sort
// a list of words. (In this case, we will be using length.)
type stringSorter struct {
	words []string
	by    By
}

// Sort allows us to sort a list of words by a given function.
func (by By) Sort(s []string) {
	sS := &stringSorter{
		words: s,
		by:    by,
	}
	sort.Sort(sS)
}

// Len is a required function for the *stringSorter interface.
func (s *stringSorter) Len() int {
	return len(s.words)
}

// Swap is a required function for the *stringSorter interface.
func (s *stringSorter) Swap(i, j int) {
	s.words[i], s.words[j] = s.words[j], s.words[i]
}

// Less is a required function for the *stringSorter interface.
func (s *stringSorter) Less(i, j int) bool {
	return s.by(s.words[i], s.words[j])
}

// sortList contains the main logic; and rather predictably,
// it sorts a list of words by how long they are.
func sortList(s []string) {
	length := func(s1, s2 string) bool {
		return len(s1) < len(s2)
	}

	By(length).Sort(s)
}
