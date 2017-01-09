package points

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/zknill/points/commands"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type slashResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

// LEADERBOARD const for the leaderboard datastore key
const LEADERBOARD string = "leaderboard"

// ENTRY const for the entry datastore key
const ENTRY string = "entry"
const token string = ""

// Run the webapp
func Run() {
	http.HandleFunc("/command", handleCommand)
}

func handleCommand(w http.ResponseWriter, r *http.Request) {
	if token != "" && r.PostFormValue("token") != token {
		http.Error(w, "Invalid Slack token.", http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	ctx := appengine.NewContext(r)
	var rtext string
	cmdStr := strings.Replace(r.PostFormValue("text"), "  ", " ", -1)
	commands := strings.Split(cmdStr, " ")

	c := commands[0]
	switch {
	case match(c, "list"):
		rtext = list(ctx)
	case match(c, "init"):
		rtext = initBoard(ctx, commands)
	case match(c, "add"):
		log.Infof(ctx, fmt.Sprintf("adding to %s, request: %s", commands[1], r))
		if strings.EqualFold(commands[1], r.Form.Get("user_name")) {
			rtext = "awww man! you cannot add points to yourself."
		} else {
			rtext = add(ctx, commands)
		}
	}

	resp := &slashResponse{
		ResponseType: "in_channel",
		Text:         rtext,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		c := appengine.NewContext(r)
		log.Errorf(c, "Error encoding JSON: %s", err)
		http.Error(w, "Error encoding JSON.", http.StatusInternalServerError)
		return
	}
}

func getResponseText(lb *points.Leaderboard) string {
	buffer := new(bytes.Buffer)
	table := points.GetTable(buffer, lb)
	table.Render()
	return "```" + buffer.String() + "```"
}

func list(ctx context.Context) (rtext string) {
	lb := &points.Leaderboard{}
	if slb, err := getLeaderboard(ctx); err == nil {
		lb.Headers = slb.Headers
	}

	if slice, err := getEntries(ctx); slice != nil {
		for _, s := range *slice {
			rtext += s.String() + " "
		}
		lb.Entries = *slice
	} else {
		return err.Error()
	}

	rtext = getResponseText(lb)
	return
}

func initBoard(ctx context.Context, commands []string) (rtext string) {
	rtext = "leaderboard exists!"
	if _, err := getLeaderboard(ctx); err != nil {
		initLeaderboard(ctx, commands[1:])
		rtext = "alright! new leaderboard"
	}
	return
}

func add(ctx context.Context, commands []string) string {
	// commands{cmd, name, num}
	if len(commands) != 2 {
		return "aww man! please use the format `/points add slackbot`"
	}

	name := commands[1]

	if entry, err := getEntry(ctx, name); err != nil {
		log.Warningf(ctx, "failed getting entry for name '%s', error: %s", name, err.Error())
	} else {
		log.Infof(ctx, "found entry using new method: %s", entry)
		entry.Points++
		storeEntry(ctx, entry)
		return fmt.Sprintf("alright! added 1 point to %s", commands[1])
	}

	return fmt.Sprintf("awwww man! something went wrong adding 1 point to %s", commands[1])
}

func match(command, matcher string) bool {
	return strings.EqualFold(command, matcher)
}
