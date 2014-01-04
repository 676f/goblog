package admin

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"fmt"
	"github.com/676f/goblog/datatypes"
	htemplate "html/template"
	"log"
	"net/http"
	ttemplate "text/template"
	"time"
)

func init() {
	http.HandleFunc("/admin/post", post)
	http.HandleFunc("/admin/post/save", save)
}

var postTemplate = htemplate.Must(htemplate.New("").ParseFiles("templates/post.html"))
var headerTemplate = ttemplate.Must(ttemplate.New("").ParseFiles("templates/header.html"))

func post(w http.ResponseWriter, r *http.Request) {
	var hpb = datatypes.WebPageBody{}

	if err := postTemplate.ExecuteTemplate(&hpb, "post.html", nil); err != nil {
		log.Fatal(err)
	}
	if err := headerTemplate.ExecuteTemplate(w, "header.html", hpb); err != nil {
		log.Fatal(err)
	}

}

func save(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := appengine.NewContext(r)

	u := user.Current(c)

	fmt.Fprint(w, u.String())
	fmt.Fprint(w, r.FormValue("title"))

	p := datatypes.Post{u.String(), r.FormValue("title"), r.FormValue("blogcontent"), time.Now()}

	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "post", nil), &p)
	if err != nil {
		log.Fatal(err)
	}
}
