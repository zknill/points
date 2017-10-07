package points

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/zknill/points/board"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"github.com/pkg/errors"
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

	msg := Parser(ctx, team, rawCommands)
	responseText := response(ctx, msg, team)

	resp := &slashResponse{
		ResponseType: "in_channel",
		Text:         responseText,
	}

	writeResponse(ctx, w, resp)
}

func response(ctx context.Context, msg Message, team string) string {
	cmd, ok := msg.(Command)
	if !ok {
		return msg.String()
	}

	standings, err := board.Load(ctx, board.NewTeam(team))
	if err != nil {
		log.Warningf(ctx, "failed to load %+s", err)
		return errorMessage(err)
	}

	resp, err := cmd.Execute(ctx, standings)
	if err != nil {
		log.Warningf(ctx, "failed to execute command %T, error: %+s", cmd, err)
		return msg.String()
	}

	return resp.String()
}

func writeResponse(ctx context.Context, w http.ResponseWriter, resp *slashResponse) {
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf(ctx, "Error encoding JSON: %s", err)
		http.Error(w, "failed!", http.StatusInternalServerError)
		return
	}
}

func errorMessage(err error) string {
	cause := errors.Cause(err)
	if msg, ok := cause.(Message); ok {
		return msg.String()
	}
	return errorText
}
