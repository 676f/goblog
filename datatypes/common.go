package datatypes

import (
	"fmt"
	"time"
)

type Post struct {
	Author, Title, Text, DateString string
	GoDate                          time.Time
	ID                              int64
}

type Posts []*Post

func (p Posts) Less(i, j int) bool {
	return int(p[i].GoDate.Sub(p[j].GoDate)) < 0
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
