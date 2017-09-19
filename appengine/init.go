package appengine

import (
	"net/http"

	"github.com/zknill/points"
)

func init() {
	http.HandleFunc("/command", points.Handler)
}
