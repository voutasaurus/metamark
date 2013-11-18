package main

import (
	"fmt"
	"testing"
)

func TestPortmanteau(t *testing.T) {
	m := map[string]bool{"common": true, "wealth": true}
	c := portmanteau("commonwealth", m)
	fmt.Println(c)
	if !c {
		t.Fail()
	}
}

func TestOtherPortmanteau(t *testing.T) {
	m := map[string]bool{"common": true, "wealth": true}
	a, b, c := portmanteau("commonly", m), portmanteau("common", m), portmanteau("wealth", m)
	fmt.Println(a, b, c)
	if a || b || c {
		t.Fail()
	}
}
