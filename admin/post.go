package admin

import (
	"fmt"
	"html/template"
	"net/http"
	"log"
	"time"
	"github.com/676f/goblog/datatypes"
	"appengine/user"
	"appengine/datastore"
	"appengine"
)

func init() {
	http.HandleFunc("/admin/post", post)
	http.HandleFunc("/admin/post/save", save)
}

func post(w http.ResponseWriter, r *http.Request) {
	var err error
	var templates *template.Template

	if templates, err = template.New("").ParseFiles("templates/post.html"); 
		err != nil { log.Fatal(err) }

	if err := templates.ExecuteTemplate(w, "post.html", nil);
		err != nil { log.Fatal(err) }

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

	u := user.Current(c); 

	fmt.Fprint(w, u.String())
	fmt.Fprint(w, r.FormValue("title"))

	p := datatypes.Post { u.String(), r.FormValue("title"), r.FormValue("blogcontent"), time.Now()}

	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "post", nil), &p)
	if err != nil {
		log.Fatal(err)
	}
}

