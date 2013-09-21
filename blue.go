package main

import (
	"fmt"
	"html/template"
	//"io/ioutil"
	"net/http"
	"regexp"
	//"strings"
	//"errors"
)

type Link struct {
	URL         string
	Description string
}

type List struct {
	Title string
	Body  []Link
}

/*
func titleDisplay(title string) string {
	return strings.Replace(title, "_", " ", -1)
}

func (p *Page) save() error {
	filename := "data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
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
	url := "http://www.wikipedia.com"
	desc := "Wikipedia: The Commie's Encyclopedia"
	link := Link{URL: url, Description: desc}
	list := []Link{link} // list := make([]Link, 1) ; append(list, link)
	ret := &List{Title: key, Body: list}

	return ret, nil
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
	if r.Method == "GET" && key != "" {
		// receive GOT data
		fmt.Print(key)
	}
	renderTemplate(w, "index", nil)
}

/*
func newHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadPage()
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "new", p)
}
*/

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
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
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
	http.HandleFunc("/", frontHandler)
	//http.HandleFunc("/index.html", indexHandler)
	// http.HandleFunc("/new/", newHandler) // Don't register yet, not implemented
	//http.HandleFunc("/old/", makeHandler(oldHandler))

	http.ListenAndServe(":8080", nil)
	fmt.Println("test")
}

/*
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
*/
