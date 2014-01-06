package goblog

import (
	"fmt"
	"github.com/676f/goblog/datatypes"
	htemplate "html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	ttemplate "text/template"

	"appengine"
	"appengine/datastore"
)

func init() {
	http.HandleFunc("/", renderPosts)
}

var siteHeaderTemplate = ttemplate.Must(ttemplate.New("").ParseFiles("templates/header.html"))
var renderTemplate = htemplate.Must(htemplate.New("").ParseFiles("templates/homepage.html", "templates/archive.html"))

func renderPosts(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if r.URL.String() == "/" {
		fmt.Fprintf(w, renderPost(c, -1, 5, "homepage.html"))
		return
	}

	splitURL := strings.SplitAfter(r.URL.String(), "/")
	switch splitURL[1] {
	case "posts/":
		postID, err := strconv.ParseInt(splitURL[2], 10, 64)
		if err != nil {
			fmt.Fprintf(w, "<p>404<br>Invalid page number with %v</p>", splitURL)
		}
		fmt.Fprintf(w, renderPost(c, postID, -1, "homepage.html"))
	case "archive": // there should be no trailing / (as in "archive/"). This should be changed. Regex "archive/?"
		fmt.Fprintf(w, renderPost(c, -1, -1, "archive.html"))
	default:
		fmt.Fprintf(w, "<p>404<br>Tried to fetch with %v</p>", splitURL)
		return
	}
}

// renderPost renders a single post if postID != -1, or it renders the specified number of posts (newest first) if numToRender != -1.
// If both parameters are -1, then it renders all posts (newest first).
// The function returns a string of the resulting templated HTML.
func renderPost(c appengine.Context, postID int64, numToRender int, postTemplateName string) string {
	hpb := datatypes.WebPageBody{}
	finalHpb := datatypes.WebPageBody{}

	q := datastore.NewQuery("post")
	if postID > -1 {
		// According to https://developers.google.com/appengine/docs/go/datastore/queries#Go_Filters,
		// Query.Filter("__key__="...) matches keys in order of
		// 1. Ancestor path
		// 2. Entity kind
		// 3 Identifier (key name or numeric ID)
		// So by making a new key with the same type and ancestor path (no ancestors) and with the same postID, the two can be matched and returned from the query.
		// There has got to be a better way to do this, and this might be a waste of a key, but since uses the same ID, maybe it doesn't matter.
		q = datastore.NewQuery("post").Filter("__key__=", datastore.NewKey(c, "post", "", postID, nil))
	} else if numToRender != -1 {
		q = datastore.NewQuery("post").Limit(numToRender)
	}

	var allPosts []*datatypes.Post
	keys, err := q.GetAll(c, &allPosts)
	if err != nil {
		return fmt.Sprintf("<p>ERROR. q.GetAll() returned `%v`</p>", err)
	}

	for i := range allPosts {
		allPosts[i].ID = keys[i].IntID()
	}
	sort.Sort(sort.Reverse(datatypes.Posts(allPosts)))

	if err := renderTemplate.ExecuteTemplate(&hpb, postTemplateName, allPosts); err != nil {
		return fmt.Sprintf("<p>ERROR. renderTemplate.Execute() returned `%v`</p>", err)
	}
	if err := siteHeaderTemplate.ExecuteTemplate(&finalHpb, "header.html", hpb); err != nil {
		return fmt.Sprintf("<p>ERROR. renderTemplate.Execute() returned `%v`</p>", err)
	}

	return finalHpb.HTMLBody
}
