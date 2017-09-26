package points

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/zknill/points/board"
	context2 "golang.org/x/net/context"
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

	standings, err := board.Load(ctx, board.NewTeam(team))
	if err != nil {
		r := &slashResponse{
			// make ephemeral response
			ResponseType: "in_channel",
		}

		var logger func(ctx context2.Context, format string, args ...interface{})
		switch errors.Cause(err).(type) {
		case board.ErrBoardNotFound:
			logger = log.Warningf
			r.Text = "create a team using `/points init`"
		default:
			logger = log.Errorf
			r.Text = "ack! something failed"
		}
		logger(ctx, "failed loading board for team: %s, error: %+s", team, err)
		writeResponse(ctx, w, r)
		return
	}

	responseText := Parser(ctx, team, rawCommands)(ctx, standings)

	resp := &slashResponse{
		ResponseType: "in_channel",
		Text:         responseText,
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
