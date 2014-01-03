package goblog

import (
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
