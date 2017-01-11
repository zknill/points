package points

import (
	"log"
	"testing"

	"github.com/zknill/points/commands"
	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
)

func Test_getResponseText(t *testing.T) {
	type args struct {
		lb *points.Leaderboard
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "simple leaderboard",
			args: args{
				&points.Leaderboard{
					Headers: []string{
						"name",
						"points",
					},
					Entries: []*points.Entry{
						{
							Name:   "Gordon",
							Points: 2,
						},
						{
							Name:   "Jane",
							Points: 1,
						},
					},
				},
			},
			want: "```+--------+--------+\n|  NAME  | POINTS |\n+--------+--------+\n| Gordon |      2 |\n| Jane   |      1 |\n+--------+--------+\n```",
		},
		{
			name: "unordered leaderboard",
			args: args{
				&points.Leaderboard{
					Headers: []string{
						"name",
						"points",
					},
					Entries: []*points.Entry{
						{
							Name:   "Jane",
							Points: 1,
						},
						{
							Name:   "Gordon",
							Points: 2,
						},
					},
				},
			},
			want: "```+--------+--------+\n|  NAME  | POINTS |\n+--------+--------+\n| Gordon |      2 |\n| Jane   |      1 |\n+--------+--------+\n```",
		},
	}
	for _, tt := range tests {
		if got := getResponseText(tt.args.lb); got != tt.want {
			t.Errorf("%q. getResponseText() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_add(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer done()

	type args struct {
		ctx      context.Context
		commands []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "simple add",
			args: args{
				ctx,
				[]string{"add", "slackbot"},
			},
			want: "alright! added a point to Slackbot",
		},
	}
	for _, tt := range tests {
		if got := add(tt.args.ctx, tt.args.commands); got != tt.want {
			t.Errorf("%q. add() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
