package goblog

import (
	"time"
)

type Post struct {
	Author, Title, Text string
	Date                time.Time
}
