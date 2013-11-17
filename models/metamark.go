package models

import (
	//"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strings"
)

const metamarksCollection = "metamarks"
const urlField = "url"

type Tag struct {
	id     string  "id"     // The tag name
	weight float64 "weight" // How relevant the tag is to the MetaMark/Search containing it
}

type MetaMark struct {
	url  string "url"  // should be Url (?)
	tags []Tag  "tags" // Metadata for this MetaMark
	talk string "talk" // Suggestion: Wikipedia style talk section (for discussing relevance of tags)
}

// CRUD OPERATIONS

// CreateMetaMark pushes the information in m to the metamark collection
// of the database as a document.
func (m *MetaMark) CreateMetaMark() error {
	return nil
}

// ReadMetaMark uses m.url to pull the matching metamark into m
// from the metamark collection of the database.
func (m *MetaMark) ReadMetaMark() error {
	return nil
}

// UpdateMetaMark uses the data in m to modify the matching metamark
// in the metamark collection of the database. Any tags not present in m
// are removed in the document corresponding to m.url.
// This function could use some thought because the changes we want to make
// will be "adding tags", "removing tags", "modifying the weight of tags",
// and from the mongo perspective we don't have to read the metamark to
// know how to update it properly.
func (m *MetaMark) UpdateMetaMark() error {
	return nil
}

// DestroyMetaMark removes the metamark in the database with url = m.url
func (m *MetaMark) DestroyMetaMark() error {
	return nil
}

// MetaSearch searches by semantic tags and returns an iterator which
// can be used to iterate over the results.
func MetaSearch(tags []Tag, metamarks *mgo.Collection) *mgo.Iter {
	return metamarks.Find(bson.M{urlField: bson.M{"$exists": true}}).Limit(10).Iter() // dummy
}

func filter(s []string, fn func(string) bool) []string {
	var p []string // == nil
	for _, v := range s {
		if fn(v) {
			p = append(p, v)
		}
	}
	return p
}

func words(s string) []string {
	nonempty := func(word string) bool {
		return word != ""
	}
	return filter(strings.Split(s, " "), nonempty)
}

// tagify takes a search term, breaks it into words, adds synonyms,
// and weights all the words with tag weights (based on whether they were
// part of the original search and how closely synonymous they are).
// The basic version of this search will add no synonyms and all weights will be
// 1 (maybe decreasing for words that are later in the search term).
// An even more advanced version would autocorrect spelling.
func tagify(s string) []Tag {
	terms := strings.Split(s, " ")
	tags := make([]Tag, 0)
	for _, term := range terms {
		if term != "" {
			// make tag and add to result
			temp := Tag{term, 1.0}
			tags = append(tags, temp)
		}
	}
	return tags
}

// SearchForPages
func SearchForPages(s string, metamarks *mgo.Collection, results chan MetaMark) {
	iter := MetaSearch(tagify(s), metamarks)
	m := MetaMark{}
	for iter.Next(&m) {
		results <- m
	}
	close(results)
	return
}
