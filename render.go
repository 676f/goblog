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
var renderTemplate = htemplate.Must(htemplate.New("").ParseFiles("templates/homepage.html", "templates/archive.html", "templates/404.html"))

func renderPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/" {
		fmt.Fprintf(w, renderPost(r, -1, 5, "homepage.html"))
		return
	}

	splitURL := strings.SplitAfter(r.URL.String(), "/")
	// if the first part of the URL ends with '/', remove the '/'
	if splitURL[1][len(splitURL[1])-1] == '/' {
		splitURL[1] = splitURL[1][:len(splitURL[1])-1]
	}

	switch splitURL[1] {
	case "posts":
		if len(splitURL) == 3 {
			postID, err := strconv.ParseInt(splitURL[2], 10, 64)
			if err != nil {
				fmt.Fprintf(w, Error404(r))
				return
			}
			fmt.Fprintf(w, renderPost(r, postID, -1, "homepage.html"))
		} else {
			fmt.Fprintf(w, renderPost(r, -1, -1, "archive.html"))
		}
	default:
		fmt.Fprintf(w, "<p>404<br>Tried to fetch with %v</p>", splitURL)
		return
	}
}

// renderPost renders a single post if postID != -1, or it renders the specified number of posts (newest first) if numToRender != -1.
// If both parameters are -1, then it renders all posts (newest first).
// The function returns a string of the resulting templated HTML.
func renderPost(r *http.Request, postID int64, numToRender int, postTemplateName string) string {
	c := appengine.NewContext(r)
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
	if len(allPosts) == 0 {
		return Error404(r)
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

func Error404(r *http.Request) string {
	hpb := datatypes.WebPageBody{}
	finalHpb := datatypes.WebPageBody{}

	if err := renderTemplate.ExecuteTemplate(&hpb, "404.html", r); err != nil {
		return fmt.Sprintf("<p>ERROR. renderTemplate.Execute() returned `%v`</p>", err)
	}
	if err := siteHeaderTemplate.ExecuteTemplate(&finalHpb, "header.html", hpb); err != nil {
		return fmt.Sprintf("<p>ERROR. renderTemplate.Execute() returned `%v`</p>", err)
	}

	return finalHpb.HTMLBody
}
