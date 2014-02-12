package render

import (
	"fmt"
	"github.com/676f/goblog/datatypes"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	ttemplate "text/template"

	"appengine"
	"appengine/datastore"
)

var siteHeaderTemplate = ttemplate.Must(ttemplate.New("").ParseFiles("templates/header.html"))
var renderTemplate = ttemplate.Must(ttemplate.New("").ParseFiles("templates/homepage.html", "templates/archive.html", "templates/404.html"))
var protocolRegexp = regexp.MustCompilePOSIX("https?://")

func RenderPosts(w http.ResponseWriter, r *http.Request) {
	url := strings.Join(protocolRegexp.Split(r.URL.String(), -1), "")
	//fmt.Fprintln(w, url)
	splitURLWithSlashes := strings.Split(url, "/")
	//fmt.Fprintln(w, splitURL, len(splitURL))

	//remove empty spaces that used to be "/"'s
	splitURL := make([]string, 1)
	for i := range splitURLWithSlashes {
		if splitURLWithSlashes[i] == "" {
			continue
		}
		splitURL = append(splitURL, splitURLWithSlashes[i])
	}
	if len(splitURLWithSlashes) < 2 {
		FinishRender(w, Error404(r))
		return
	}

	// if there is no path after domain name
	if r.URL.String() == "/" {
		q, err := getQuery(r, -1, 5, "posts")
		if err != nil {
			FinishRender(w, fmt.Sprintf("<p>%v</p>", err))
			return
		}
		FinishRender(w, renderPost(q, "homepage.html"))
		return
	}

	// no type means get all from every type
	queryType := "posts"
	q := datatypes.Posts{}
	var postID int64 = -1
	numPostsToGet := -1
	templatePage := ""
	var err error

	switch splitURL[1] {
	case "posts":
		// if a post ID is attached
		if len(splitURL) == 3 {
			postID, err = strconv.ParseInt(splitURL[2], 10, 64)
			if err != nil {
				FinishRender(w, Error404(r))
				return
			}
		}
	case "tags":
		switch len(splitURL) {
		case 3:
			queryType = splitURL[2]
		case 4:
			// if a post ID is attached
			postID, err = strconv.ParseInt(splitURL[2], 10, 64)
			if err != nil {
				FinishRender(w, Error404(r))
				return
			}
		}
	}

	if postID != -1 {
		// then render only that post and none else
		templatePage = "homepage.html"
		numPostsToGet = 1
	} else {
		// if there is no post ID
		// then we want to render the archive: all posts
		// the numPostsToGet is -1 (all)
		templatePage = "archive.html"
	}

	q, err = getQuery(r, postID, numPostsToGet, queryType)
	if err != nil {
		FinishRender(w, fmt.Sprintf("<p>%v</p>", err))
		return
	}

	FinishRender(w, renderPost(q, templatePage))
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

func getQuery(r *http.Request, postID int64, numToRender int, queryType string) ([]*datatypes.Post, error) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery(queryType)
	if postID > -1 {
		// According to https://developers.google.com/appengine/docs/go/datastore/queries#Go_Filters,
		// Query.Filter("__key__="...) matches keys in order of
		// 1. Ancestor path
		// 2. Entity kind
		// 3 Identifier (key name or numeric ID)
		// So by making a new key with the same type and ancestor path (no ancestors) and with the same postID, the two can be matched and returned from the query.
		// There has got to be a better way to do this, and this might be a waste of a key, but since uses the same ID, maybe it doesn't matter.
		q = datastore.NewQuery(queryType).Filter("__key__=", datastore.NewKey(c, queryType, "", postID, nil)).Order("-GoDate")
	} else if numToRender != -1 {
		q = datastore.NewQuery(queryType).Order("-GoDate").Limit(numToRender)
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
