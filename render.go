package goblog

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	//"time"

	"appengine"
	"appengine/datastore"
)

const renderTemplateHTML = `
<div class="post">
	<div class="post_header">written by <b>{{.Author}}</b> on <b>{{.Date}}</b></div>
	<br/>
	<div class="post_text">{{.Text}}</div>
</div>
`

const siteHeader = `
<html>
	<head>
		<title>Goblog Testing Site</title>	
		<link rel="stylesheet" type="text/css" href="/stylesheets/main.css"/>
	</head>
	<body>
		<div id="top_bar">
			<div id="header"><a href="/">Goblog Testing Site</a></div>
			<div id="header_accent"></div>
		</div>
	<div id="main">
		<div id="sidebar">
			This is a testing site for the sick Goblog blogging.
		</div>
	{{range .}}
	<div class="post">
	       <div class="post_header">written by <b>{{.Author}}</b> on <b>{{.Date}}</b></div>
	       <br/>
	       <div class="post_text">{{.Text}}</div>
	</div>
	{{end}}
	</div>
	</body>
</html>
`

var renderTemplate = template.Must(template.New("post").Parse(renderTemplateHTML))
var siteHeaderTemplate = template.Must(template.New("site").Parse(siteHeader))

func RenderAllPosts(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	//fmt.Fprintf(w, siteHeader)

	//p1 := Post{
	//"Joe Smith",
	//"Title1",
	//"This is an awesome blog.",
	//time.Now(),
	//}
	//fmt.Fprintf(w, "<p>%v</p>", p1)

	//_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "post", nil), &p1)
	//if err != nil {
	//fmt.Fprintf(w, "<p>ERROR. datastore.Put() returned `%v`</p>", err)
	//}

	q := datastore.NewQuery("post")
	//count, err := q.Count(c)
	//if err != nil {
	//fmt.Fprintf(w, "<p>ERROR. q.Count() returned '%v'</p>", err)
	//}
	//fmt.Fprintf(w, "<p>%v</p>", count)

	var allPosts Posts
	_, err := q.GetAll(c, &allPosts)
	if err != nil {
		fmt.Fprintf(w, "<p>ERROR. q.GetAll() returned `%v`</p>", err)
	}

	sort.Sort(sort.Reverse(allPosts))
	if err := siteHeaderTemplate.Execute(w, []*Post(allPosts)); err != nil {
		fmt.Fprintf(w, "<p>ERROR. renderTemplate.Execute() returned `%v`</p>", err)
	}
}
