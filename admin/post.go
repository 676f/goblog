package admin

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"github.com/676f/goblog/datatypes"
	htemplate "html/template"
	"log"
	"net/http"
	"strconv"
	ttemplate "text/template"
	"time"
)

func init() {
	http.HandleFunc("/admin/post", post)
}

func post(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		save(w, r)
	} else {

		var postTemplate = htemplate.Must(htemplate.New("").ParseFiles("templates/post.html"))
		var headerTemplate = ttemplate.Must(ttemplate.New("").ParseFiles("templates/header.html"))
		var hpb = datatypes.WebPageBody{}

		if err := postTemplate.ExecuteTemplate(&hpb, "post.html", nil); err != nil {
			log.Fatal(err)
		}
		if err := headerTemplate.ExecuteTemplate(w, "header.html", hpb); err != nil {
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

	p := datatypes.Post{u.String(), r.FormValue("title"), r.FormValue("blogcontent"), time.Now(), -1}
	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "post", nil), &p)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/posts/"+strconv.FormatInt(key.IntID(), 10), http.StatusFound)
}
