package main

import (
		"fmt"
		"html/template"
		"io/ioutil"
		"net/http"
		//"regexp"
		"strings"
		//"errors"
)

type Page struct {
	Title string
	Body []byte
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

*/

func frontHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/Front_Page", http.StatusFound)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadPage()
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "new", p)
}

/* // Activate for functions on templates
var funcMap = template.FuncMap{
			"titleFmt":titleDisplay,
	}

var templates = template.Must(template.New("titleTest").Funcs(funcMap).ParseFiles("tmpl/edit.html", "tmpl/view.html"))
*/

var templates = template.Must(template.ParseFiles("tmpl/main.html", "tmpl/newlist.html"))

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


func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		fmt.Println("Template " + tmpl + " cannot be rendered")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func main() {
	http.HandleFunc("/", frontHandler)
  
  	http.HandleFunc("/new/", newHandler)
  
	http.ListenAndServe(":8080", nil)
	fmt.Println("test")
}

/*	
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
*/

