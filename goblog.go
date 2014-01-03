package goblog

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", RenderAllPosts)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Heasfasdllo, world!")
}
