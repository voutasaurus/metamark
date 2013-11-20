// Package handlers contains the handlers for the web server
// portion of Blue.
package handlers

import (
	"fmt"
	"github.com/voutasaurus/Blue/models"
	"html/template"
	"net/http"
)

// GetList is a channel for handling requests to retrieve lists.
var GetList = make(chan models.ListRetrieve)

// AddList is a channel for handling requests to add new lists.
var AddList = make(chan models.AddRequest)

// RemoveList is a channel for handling requests to remove lists.
var RemoveList = make(chan string)

// lastList is a hacky solution to displaying the previously
// accessed list.
var lastList InfoBookmarks

// templates is a list of HTML templates which are parsed for the
// handler.
var templates = template.Must(template.ParseFiles(
	"www/index.html",
	"www/new.html",
	"www/old.html"))

// InfoBookmarks is a hacky solution type which contains a list
// of bookmarks, along with a message for the HTML template
// to display (such as an error, or as the type is named, information.)
type InfoBookmarks struct {
	List    models.Bookmarks
	Message string
}

// MakeRedirHandler creates and serves a handler for redirecting pages.
// It is possibly misnamed; it only ever redirects users back to the
// index page at the "/" path.
func MakeRedirHandler(pass string) func(http.ResponseWriter, *http.Request) {
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

// NewHandler creates and serves a handler for the "create a new list"
// page.
func NewHandler(w http.ResponseWriter, r *http.Request) {
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
				//fmt.Println(urls, descs)
				bookmarks := new(models.Bookmarks)
				for key := range urls {
					bookmarks.List = append(bookmarks.List, models.Bookmark{"http://" + urls[key], descs[key]})
					//fmt.Println(bookmarks.List)
				}
				reply := make(chan string)
				AddList <- models.AddRequest{*bookmarks, reply}
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
		MakeRedirHandler(r.FormValue("pass"))(w, r)
	default: // we just got here! come on guys, seriously...
		renderTemplate(w, "new", &InfoBookmarks{Message: "There's nothing here yet..."})
	}
}

// loadList retrieves an existing list from the database.
func loadList(key string) *InfoBookmarks {
	// implement later
	reply := make(chan models.Bookmarks)
	GetList <- models.ListRetrieve{key, reply}
	newList := <-reply
	return &InfoBookmarks{newList, ""}
}

// renderTemplate displays the correctly templated HTML so the
// user can access the page.
func renderTemplate(w http.ResponseWriter, tmpl string, l *InfoBookmarks) {
	err := templates.ExecuteTemplate(w, tmpl+".html", l)
	if err != nil {
		fmt.Println("Template " + tmpl + " cannot be rendered")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
