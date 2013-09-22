package main

import (
	"fmt"
	"html/template"
	//"io/ioutil"
	"net/http"
	//"regexp"
	//"strings"
	//"errors"
	"github.com/voutasaurus/Blue/Models"
)

// var urlLists = make(map[string]List)

var getList = make(chan models.ListRetrieve)
var addList = make(chan models.AddRequest)
var removeList = make(chan string)
var lastList InfoBookmarks

type InfoBookmarks struct {
	List    models.Bookmarks
	Message string
}

func loadList(key string) *InfoBookmarks {
	// implement later
	reply := make(chan models.Bookmarks)
	getList <- models.ListRetrieve{key, reply}
	newList := <-reply
	return &InfoBookmarks{newList, ""}
}

func makeRedirHandler(pass string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var key string
		if pass != "" {
			key = pass
		} else {
			key = r.FormValue("pass")
		}
		if r.Method == "POST" && key != "" {
			// receive POSTed data
			l := loadList(key)
			if l.List.Key == "" {
				l.Message = "We couldn't find a list with your key \"" + key + "\"."
			}
			lastList = *l
			renderTemplate(w, "index", l)
		} else {
			renderTemplate(w, "index", new(InfoBookmarks))
		}
	}
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	caller := r.FormValue("formId")
	fmt.Println(caller)
	switch caller {
	case "addUrl": // do this once db is up and running
		if r.Method == "POST" {
			l := lastList
			err := r.ParseForm()
			if err == nil { // this should never fail
				urls := r.Form["url"]
				descs := r.Form["desc"]
				fmt.Println(urls, descs)
				bookmarks := new(models.Bookmarks)
				for key := range urls {
					bookmarks.List = append(bookmarks.List, models.Bookmark{"http://" + urls[key], descs[key]})
					fmt.Println(bookmarks.List)
				}
				reply := make(chan string)
				addList <- models.AddRequest{*bookmarks, reply}
				newKey := <-reply
				if newKey != "" {
					l = *loadList(newKey)
					l.Message = "Successfully saved your links with key \"" + newKey + "\"."
				}
			}
			lastList = l
			renderTemplate(w, "new", &l)
		}
	case "getList": // wrong page, try ./
		makeRedirHandler(r.FormValue("pass"))(w, r)
	default: // we just got here! come on guys, seriously...
		renderTemplate(w, "new", &InfoBookmarks{Message: "There's nothing here yet..."})
	}

}

var templates = template.Must(template.ParseFiles("www/index.html", "www/new.html", "www/old.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, l *InfoBookmarks) {
	err := templates.ExecuteTemplate(w, tmpl+".html", l)
	if err != nil {
		fmt.Println("Template " + tmpl + " cannot be rendered")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	go models.BookmarksCollection(getList, addList, removeList)

	http.HandleFunc("/new", newHandler)
	http.HandleFunc("/", makeRedirHandler(""))
	http.Handle("/javascripts/", http.FileServer(http.Dir("www")))
	http.Handle("/stylesheets/", http.FileServer(http.Dir("www")))

	http.ListenAndServe(":8080", nil)
	fmt.Println("test")
}
