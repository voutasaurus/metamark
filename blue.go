// Package main is the main package of the metamark library.
//
// It initializes important channels, and controls the
// fileserver handling the main pages.
package main

import (
	"fmt"
	//"html/template"
	//"io/ioutil"
	"net/http"
	//"regexp"
	//"strings"
	//"errors"
	"github.com/voutasaurus/metamark/handlers"
	"github.com/voutasaurus/metamark/models"
)

// var urlLists = make(map[string]List)

func main() {
	defer close(handlers.GetList)
	defer close(handlers.AddList)
	defer close(handlers.RemoveList)
	go models.BookmarksCollection(handlers.GetList, handlers.AddList, handlers.RemoveList)

	http.HandleFunc("/new", handlers.NewHandler)
	http.HandleFunc("/", handlers.MakeRedirHandler(""))
	http.Handle("/javascripts/", http.FileServer(http.Dir("www")))
	http.Handle("/stylesheets/", http.FileServer(http.Dir("www")))

	fmt.Println("metamark server up and running...") // does main ever get here? - It does now. :P
	http.ListenAndServe(":8080", nil)
}
