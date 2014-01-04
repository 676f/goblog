package admin

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"github.com/676f/goblog/datatypes"
	"html/template"
	"log"
	"net/http"
	"time"
)

func init() {
	http.HandleFunc("/admin/post", post)
}

func post(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		save(w, r)
	} else {
		var err error
		var templates *template.Template

		if templates, err = template.New("").ParseFiles("templates/post.html"); err != nil {
			log.Fatal(err)
		}

		if err := templates.ExecuteTemplate(w, "post.html", nil); err != nil {
			log.Fatal(err)
		}
	}

}

func save(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := appengine.NewContext(r)

	u := user.Current(c)

	p := datatypes.Post{u.String(), r.FormValue("title"), r.FormValue("blogcontent"), time.Now()}

	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "post", nil), &p)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
