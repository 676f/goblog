package render

import (
	"fmt"
	"github.com/676f/goblog/datatypes"
//	htemplate "html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	ttemplate "text/template"

	"appengine"
	"appengine/datastore"
)

var siteHeaderTemplate = ttemplate.Must(ttemplate.New("").ParseFiles("templates/header.html"))
var renderTemplate = ttemplate.Must(ttemplate.New("").ParseFiles("templates/homepage.html", "templates/archive.html", "templates/404.html"))

func RenderPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/" {
		q, err := getQuery(r, -1, 5)
		if err != nil {
			FinishRender(w, fmt.Sprintf("<p>%v</p>", err))
			return
		}
		FinishRender(w, renderPost(q, "homepage.html"))
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
				FinishRender(w, Error404(r))
				return
			}
			q, err := getQuery(r, postID, 1)
			if len(q) <= 0 {
				FinishRender(w, Error404(r))
				return
			} else if err != nil {
				FinishRender(w, fmt.Sprintf("<p>%v</p>", err))
				return
			}
			FinishRender(w, renderPost(q, "homepage.html"))
		} else {
			q, err := getQuery(r, -1, -1)
			if err != nil {
				FinishRender(w, fmt.Sprintf("<p>%v</p>", err))
			}
			FinishRender(w, renderPost(q, "archive.html"))
		}
	default:
		FinishRender(w, Error404(r))
		return
	}
}

// renderPost renders a single post if postID != -1, or it renders the specified number of posts (newest first) if numToRender != -1.
// If both parameters are -1, then it renders all posts (newest first).
// The function returns a string of the resulting templated HTML.
func renderPost(allPosts []*datatypes.Post, postTemplateName string) string {
	hpb := datatypes.WebPageBody{}
	if err := renderTemplate.ExecuteTemplate(&hpb, postTemplateName, allPosts); err != nil {
		return fmt.Sprintf("<p>ERROR. renderTemplate.Execute() returned `%v`</p>", err)
	}

	return hpb.HTMLBody
}

func getQuery(r *http.Request, postID int64, numToRender int) ([]*datatypes.Post, error) {
	c := appengine.NewContext(r)
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
		q = datastore.NewQuery("post").Order("-Date").Limit(numToRender)
	}

	var allPosts []*datatypes.Post
	keys, err := q.GetAll(c, &allPosts)
	if err != nil {
		return nil, err
	}

	for i := range allPosts {
		allPosts[i].ID = keys[i].IntID()
	}
	sort.Sort(sort.Reverse(datatypes.Posts(allPosts)))

	return allPosts, nil
}

func Error404(r *http.Request) string {
	hpb := datatypes.WebPageBody{}
	if err := renderTemplate.ExecuteTemplate(&hpb, "404.html", r); err != nil {
		return fmt.Sprintf("<p>ERROR. renderTemplate.Execute() returned `%v`</p>", err)
	}

	return hpb.HTMLBody
}

func FinishRender(w http.ResponseWriter, hpb string) {
	finalHpb := datatypes.WebPageBody{}
	if err := siteHeaderTemplate.ExecuteTemplate(&finalHpb, "header.html", hpb); err != nil {
		fmt.Fprintf(w, "<p>ERROR. renderTemplate.Execute() returned `%v`</p>", err)
		return
	}

	fmt.Fprint(w, finalHpb.HTMLBody)
}
