package datatypes

import (
	"fmt"
	"time"
)

type Post struct {
	Author, Title, Text string
	Date                time.Time
}

type Posts []*Post

func (p Posts) Less(i, j int) bool {
	return int(p[i].Date.Sub(p[j].Date)) < 0
}

func (p Posts) Len() int { return len(p) }

func (p Posts) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type WebPageBody struct {
	HTMLBody string
}

func (wpb *WebPageBody) Write(p []byte) (n int, err error) {
	//wpb.HTMLBody = fmt.Sprint(len(p))
	for i := range p {
		wpb.HTMLBody = fmt.Sprintf("%s%c", wpb.HTMLBody, p[i])
	}
	return len(p), nil
}
