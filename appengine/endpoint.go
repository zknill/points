package appengine

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/zknill/points"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type slashResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

// Handler implements the http.Handler interface and handles the requests made by slack users.
func Handler(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("SLACK_TOKEN")
	if token != "" && r.PostFormValue("token") != token {
		http.Error(w, "Invalid Slack token.", http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	ctx := appengine.NewContext(r)

	rawCommands := r.PostFormValue("text")
	team := r.PostFormValue("team_domain")

	msg := points.Controller.Parse(ctx, team, rawCommands)

	resp := &slashResponse{
		ResponseType: "in_channel",
		Text:         msg,
	}

	writeResponse(ctx, w, resp)
}

func writeResponse(ctx context.Context, w http.ResponseWriter, resp *slashResponse) {
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf(ctx, "Error encoding JSON: %s", err)
		http.Error(w, "failed!", http.StatusInternalServerError)
		return
	}
}
