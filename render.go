package goblog

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"

	"appengine"
	"appengine/datastore"
)

var siteHeaderTemplate = template.Must(template.New("").ParseFiles("templates/homepage.html"))

func RenderAllPosts(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	q := datastore.NewQuery("post")
	var allPosts Posts
	_, err := q.GetAll(c, &allPosts)
	if err != nil {
		fmt.Fprintf(w, "<p>ERROR. q.GetAll() returned `%v`</p>", err)
	}

	sort.Sort(sort.Reverse(allPosts))
	if err := siteHeaderTemplate.ExecuteTemplate(w, "homepage.html", []*Post(allPosts)); err != nil {
		fmt.Fprintf(w, "<p>ERROR. renderTemplate.Execute() returned `%v`</p>", err)
	}
}
