package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

//Cache Template
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

type Page struct {
	Title string
	Body  []byte
}

func (p Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	if body, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return &Page{Title: title, Body: body}, nil
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	page, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	renderTemplate(w, "view", page)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}
	renderTemplate(w, "edit", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	if err := p.save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	/*
	 *t, err := template.ParseFiles(tmpl + ".html")
	 *if err != nil {
	 *  http.Error(w, err.Error(), http.StatusInternalServerError)
	 *  return
	 *}
	 t.Execute(w,p)
	*/

	if err := templates.ExecuteTemplate(w, tmpl+".html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
	/*
	 *p1 := Page{Title: "Test", Body: []byte("This is test page")}
	 *p1.save()
	 *p2, _ := loadPage("Test")
	 *fmt.Println(string(p2.Body))
	 */
}
