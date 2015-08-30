package models

import (
  	"testing"
  	"fmt"
  	"strings"
)

func TestSplit(t *testing.T) {
  	search := "Hello my name is Frank"
  	terms := strings.Split(search, " ")
    for _, term := range terms {
        fmt.Println(term)
    }
    fmt.Println(len(terms))
}

func TestTagify(t *testing.T) {
  	search := "Hello my name is Frank"
  	tags := tagify(search)
  	fmt.Println(tags)
}