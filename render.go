package goblog

import (
	"fmt"
	"github.com/676f/goblog/datatypes"
	htemplate "html/template"
	"net/http"
	"sort"
	ttemplate "text/template"

	"appengine"
	"appengine/datastore"
)

var siteHeaderTemplate = ttemplate.Must(ttemplate.New("").ParseFiles("templates/header.html"))
var homePageBodyTemplate = htemplate.Must(htemplate.New("").ParseFiles("templates/homepage.html"))

func RenderAllPosts(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	hpb := datatypes.WebPageBody{}

	q := datastore.NewQuery("post")
	var allPosts []*datatypes.Post
	_, err := q.GetAll(c, &allPosts)
	if err != nil {
		fmt.Fprintf(w, "<p>ERROR. q.GetAll() returned `%v`</p>", err)
	}

	sort.Sort(sort.Reverse(datatypes.Posts(allPosts)))
	if err := homePageBodyTemplate.ExecuteTemplate(&hpb, "homepage.html", allPosts); err != nil {
		fmt.Fprintf(w, "<p>ERROR. renderTemplate.Execute() returned `%v`</p>", err)
	}
	if err := siteHeaderTemplate.ExecuteTemplate(w, "header.html", hpb); err != nil {
		fmt.Fprintf(w, "<p>ERROR. renderTemplate.Execute() returned `%v`</p>", err)
	}
}
