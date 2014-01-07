package goblog

import (
	"github.com/676f/goblog/admin"
	"github.com/676f/goblog/render"
	"net/http"
)

func init() {
	http.HandleFunc("/", render.RenderPosts)
	http.HandleFunc("/admin/post", admin.Post)
}
