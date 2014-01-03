package goblog

import (
	"fmt"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
)

func RenderAllPosts(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	p1 := Post{
		"Joe",
		"Title1",
		"This is an awesome blog.",
		time.Now(),
	}
	fmt.Fprintln(w, p1)

	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "post", nil), &p1)
	if err != nil {
		fmt.Fprintf(w, "ERROR. datastore.Put() returned `%v`\n", err)
	}

	q := datastore.NewQuery("post")
	count, err := q.Count(c)
	if err != nil {
		fmt.Fprintf(w, "ERROR. q.Count() returned '%v'\n", err)
	}
	fmt.Fprintln(w, count)

	var allPosts []*Post
	_, err = q.GetAll(c, &allPosts)
	if err != nil {
		fmt.Fprintf(w, "ERROR. q.GetAll() returned `%v`\n", err)
	}

	for i := range allPosts {
		fmt.Fprintln(w, *allPosts[i])
	}
}
