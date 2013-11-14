package main

import (
	"fmt"
	//"html/template"
	//"io/ioutil"
	"net/http"
	//"regexp"
	//"strings"
	//"errors"
	"github.com/voutasaurus/Blue/handlers"
	"github.com/voutasaurus/Blue/models"
)

// var urlLists = make(map[string]List)

func main() {
	go models.BookmarksCollection(handlers.GetList, handlers.AddList, handlers.RemoveList)

	http.HandleFunc("/new", handlers.NewHandler)
	http.HandleFunc("/", handlers.MakeRedirHandler(""))
	http.Handle("/javascripts/", http.FileServer(http.Dir("www")))
	http.Handle("/stylesheets/", http.FileServer(http.Dir("www")))

	http.ListenAndServe(":8080", nil)
	fmt.Println("Project Blue server up and running...") // does main ever get here?
}
