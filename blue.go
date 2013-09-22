package main

import (
	"fmt"
	"html/template"
	//"io/ioutil"
	"net/http"
	"regexp"
	//"strings"
	"errors"
	"github.com/voutasaurus/Blue/Models"
)

var hack *models.Bookmark

var urlLists = make(map[string]List)

var getList = make(chan models.ListRetrieve)
var addList = make(chan models.AddRequest)
var removeList = make(chan string)

type Link struct {
	URL         string
	Description string
}

type List struct {
	Title string
	Links []Link
	Error string
}

func init() {
	// Initialize a dummy "urllists"
	url1 := "http://www.wikipedia.com"
	url2 := "http://www.wiktionary.org"
	desc1 := "Wikipedia: The Commie's Encyclopedia"
	desc2 := "Wiktionary: Because Reds Need Definitions Too"
	link1 := Link{URL: url1, Description: desc1}
	link2 := Link{URL: url2, Description: desc2}
	list := []Link{link1, link2} // list := make([]Link, 1) ; append(list, link)
	urlLists["Wiki"] = List{Title: "Wiki", Links: list}

	url3 := "http://www.youtube.com"
	url4 := "http://www.dailymotion.com"
	url5 := "http://www.vimeo.com"
	desc3 := "Fuck videos"
	desc4 := "Videos are shit"
	desc5 := "What the fuck is this"
	link3 := Link{URL: url3, Description: desc3}
	link4 := Link{URL: url4, Description: desc4}
	link5 := Link{URL: url5, Description: desc5}
	list = []Link{link3, link4, link5} // list := make([]Link, 1) ; append(list, link)
	urlLists["Vids"] = List{Title: "Vids", Links: list}
}

/*
func titleDisplay(title string) string {
	return strings.Replace(title, "_", " ", -1)
}

func (p *Page) save() error {
	filename := "data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Links, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	Links, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Links: Links}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}
*/

func loadList(key string) (*List, error) {
	// implement later
	retVal, ok := urlLists[key]
	if ok {
		return &retVal, nil
	} else {
		return nil, errors.New("No such key \"" + key + "\".")
	}
}

const lenPath = len("/?pass=")

var titleValidator = regexp.MustCompile("^[a-zA-Z\\+]+$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[lenPath:]
		if !titleValidator.MatchString(title) {
			http.NotFound(w, r)
			return //err = errors.New("Invalid Page Title")
		}
		fn(w, r, title)
	}
}

func frontHandler(w http.ResponseWriter, r *http.Request) {
	//http.Redirect(w, r, "/index.html", http.StatusFound)
	key := r.FormValue("pass")
	if r.Method == "POST" && key != "" {
		// receive POSTed data
		l, err := loadList(key)
		if err != nil {
			l = &List{Error: err.Error()}
		}
		renderTemplate(w, "index", l)
	} else {
		renderTemplate(w, "index", new(List))
	}
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
			l, err := loadList(key)
			if err != nil {
				l = &List{Error: err.Error()}
			}
			renderTemplate(w, "index", l)
		} else {
			renderTemplate(w, "index", new(List))
		}
	}
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	caller := r.FormValue("formId")
	switch caller {
	case "newUrl": // do this once db is up and running
		key := r.FormValue("pass")
		if r.Method == "POST" && key != "" {
			l, err := loadList(key)
			if err != nil {
				l = &List{Error: err.Error()}
			} else {
				// receive POSTed data
				l.Links = append(l.Links, Link{r.FormValue("url"), r.FormValue("desc")})
				l.Error = "Successfully added \"" + r.FormValue("desc") + "\"."
				urlLists[key] = *l
			}
			renderTemplate(w, "new", l)
		}
	case "getList": // wrong page, try ./
		makeRedirHandler(r.FormValue("pass"))(w, r)
	default: // we just got here! come on guys, seriously...
		renderTemplate(w, "new", new(List))
	}

}

func oldHandler(w http.ResponseWriter, r *http.Request, key string) {
	l, err := loadList(key)
	if err != nil {
		l = &List{}
	}
	renderTemplate(w, "old", l)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

/* // Activate for functions on templates
var funcMap = template.FuncMap{
			"titleFmt":titleDisplay,
	}

var templates = template.Must(template.New("titleTest").Funcs(funcMap).ParseFiles("tmpl/edit.html", "tmpl/view.html"))
*/

var templates = template.Must(template.ParseFiles("www/index.html", "www/new.html", "www/old.html"))

/* // Parts may be salvagable for newlist form action
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	Links := r.FormValue("Links")
	p := &Page{Title: title, Links: []byte(Links)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound) // Should redirect to list retrieve
}

*/

func renderTemplate(w http.ResponseWriter, tmpl string, l *List) {
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
	//http.HandleFunc("/index.html", indexHandler)
	// http.HandleFunc("/new/", newHandler) // Don't register yet, not implemented
	//http.HandleFunc("/old/", makeHandler(oldHandler))

	http.ListenAndServe(":8080", nil)
	fmt.Println("test")
}

/*
	p1 := &Page{Title: "TestPage", Links: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Links))
*/
