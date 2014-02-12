package admin

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"github.com/676f/goblog/datatypes"
	"github.com/676f/goblog/render"
	htemplate "html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Post(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		save(w, r)
	} else {

		var postTemplate = htemplate.Must(htemplate.New("").ParseFiles("templates/post.html"))
		var hpb = datatypes.WebPageBody{}

		if err := postTemplate.ExecuteTemplate(&hpb, "post.html", nil); err != nil {
			log.Fatal(err)
		}
		render.FinishRender(w, hpb.HTMLBody)
	}

}

func save(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := appengine.NewContext(r)
	u := user.Current(c)

	tags := strings.Split(r.FormValue("tags"), ",")
	p := datatypes.Post{
		Author: u.String(),
		Title:  r.FormValue("title"),
		Text:   r.FormValue("blogcontent"),
		GoDate: time.Now(),
		ID:     -1,
		Tags:   tags,
	}
	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "posts", nil), &p)
	if err != nil {
		log.Fatal(err)
	}
	for i := range tags {
		_, err := datastore.Put(c, datastore.NewIncompleteKey(c, tags[i], nil), &p)
		if err != nil {
			log.Fatal(err)
		}
	}

	time.Sleep(500 * time.Millisecond)
	http.Redirect(w, r, "/posts/"+strconv.FormatInt(key.IntID(), 10), http.StatusFound)
}
