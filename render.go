package goblog

import (
	"fmt"
	"github.com/676f/goblog/datatypes"
	htemplate "html/template"
	"net/http"
	"sort"
	"strings"
	ttemplate "text/template"

	"appengine"
	"appengine/datastore"
)

func init() {
	http.HandleFunc("/", renderPosts)
}

var siteHeaderTemplate = ttemplate.Must(ttemplate.New("").ParseFiles("templates/header.html"))
var homePageBodyTemplate = htemplate.Must(htemplate.New("").ParseFiles("templates/homepage.html"))

func renderPosts(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	splitURL := strings.SplitAfter(r.URL.String(), "/")
	// the split slice ends with a blank entry
	postNum := splitURL[len(splitURL)-2]
	if postNum != "/" {
		fmt.Fprintf(w, "other url: %v\n", splitURL)
		return
	}

	fmt.Fprintf(w, renderPost(c, -1, 3))
}

// renderPost renders a single post if postID != -1, or it renders the specified number of posts (newest first) if numToRender != -1.
// If both parameters are -1, then it renders all posts (newest first).
// The function returns a string of the resulting templated HTML.
func renderPost(c appengine.Context, postID, numToRender int) string {
	hpb := datatypes.WebPageBody{}
	finalHpb := datatypes.WebPageBody{}

	q := datastore.NewQuery("post")
	var allPosts []*datatypes.Post
	_, err := q.GetAll(c, &allPosts)
	if err != nil {
		return fmt.Sprintf("<p>ERROR. q.GetAll() returned `%v`</p>", err)
	}
	sort.Sort(sort.Reverse(datatypes.Posts(allPosts)))

	if len(allPosts) > numToRender && numToRender != -1 {
		allPosts = allPosts[:numToRender]
	}

	if err := homePageBodyTemplate.ExecuteTemplate(&hpb, "homepage.html", allPosts); err != nil {
		return fmt.Sprintf("<p>ERROR. renderTemplate.Execute() returned `%v`</p>", err)
	}
	if err := siteHeaderTemplate.ExecuteTemplate(&finalHpb, "header.html", hpb); err != nil {
		return fmt.Sprintf("<p>ERROR. renderTemplate.Execute() returned `%v`</p>", err)
	}

	return finalHpb.HTMLBody
}
