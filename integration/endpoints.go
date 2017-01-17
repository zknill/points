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

type operation func(ctx context.Context, user string, commands ...string) string

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

	cmdStr := strings.Replace(r.PostFormValue("text"), "  ", " ", -1)
	commands := strings.Split(cmdStr, " ")
	user := r.Form.Get("user_name")

	rtext := matcher(commands[0])(ctx, user, commands...)

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

func list(ctx context.Context, _ string, _ ...string) (rtext string) {
	lb := &points.Leaderboard{}
	if slb, err := getLeaderboard(ctx); err == nil {
		lb.Headers = slb.Headers
	}

	if slice, err := getEntries(ctx); slice != nil {
		lb.Entries = *slice
	} else {
		log.Warningf(ctx, "something went wrong getting entries, error: %s", err.Error())
		return "awww man! something went wrong..."
	}

	rtext = getResponseText(lb)
	return
}

func initBoard(ctx context.Context, _ string, commands ...string) (rtext string) {
	rtext = "leaderboard exists!"
	if _, err := getLeaderboard(ctx); err != nil {
		if initErr := initLeaderboard(ctx, commands[1:]); err != nil {
			log.Warningf(ctx, "failed to init leaderboard, error: %s", initErr.Error())
			return "awww man! something went wrong setting up your leaderboard..."
		}
		rtext = "alright! new leaderboard"
	}
	return
}

func add(ctx context.Context, user string, commands ...string) string {
	// commands{cmd, name, num}
	if len(commands) != 2 {
		return "aww man! please use the format `/points add slackbot`"
	}

	if strings.EqualFold(commands[1], user) {
		return fmt.Sprintf("sorry %s! you cannot add points to yourself", user)
	}

	name := commands[1]
	entry, err := getEntry(ctx, name)
	if err != nil {
		log.Infof(ctx, "failed getting entry for name '%s', error: %s", name, err.Error())
		entry = &points.Entry{Name: strings.Title(name), Points: 1}
	} else {
		log.Infof(ctx, "found entry using new method: %s", entry)
		entry.Points++
	}
	if err := storeEntry(ctx, entry); err != nil {
		return fmt.Sprintf("awww man! something went wrong adding a point to %s", strings.Title(name))
	}
	return fmt.Sprintf("alright! added a point to %s", strings.Title(commands[1]))
}

func unknown(ctx context.Context, user string, commands ...string) string {
	log.Infof(ctx, "%s not a recognised command", commands[0])
	return fmt.Sprintf("awww man! sorry %s but %s is not a supported command", user, commands[0])
}

func matcher(c string) operation {
	switch {
	case strings.EqualFold("list", c):
		return list
	case strings.EqualFold("add", c):
		return add
	case strings.EqualFold("init", c):
		return initBoard
	default:
		return unknown
	}
}
