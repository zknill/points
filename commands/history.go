package points

import (
	"strings"
	"time"
)

// History fir the leaderboard
type History struct {
	Timestamp int64
	Message   string
	Args      []string
}

func (lb *Leaderboard) addHistory(command string, args ...string) {
	lb.History = append(lb.History, &History{time.Now().UnixNano(), command, args})
}

// String formats the history for printing
func (h *History) String() string {
	t := time.Unix(0, h.Timestamp)
	timeStr := t.Format("2006-01-02 15:04:05")
	return timeStr + ": " + h.Message + " " + strings.Join(h.Args, " ")
}
