package points

import (
	"flag"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/urfave/cli"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		file string
		cmds []string
		want int
	}{
		{
			name: "simple add",
			file: "./test-fixtures/add.json",
			cmds: []string{"slackbot", "2"},
			want: 3,
		},
		{
			name: "only name",
			file: "./test-fixtures/add.json",
			cmds: []string{"slackbot"},
			want: 2,
		},
	}
	for _, tt := range tests {
		ctx := newApp(tt.file, tt.cmds)
		lbSetup(tt.file)
		Add(ctx)
		got := currPoints(Read(ctx), tt.cmds[0])
		if tt.want != got {
			t.Errorf("%q. Add() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestReset(t *testing.T) {
	tests := []struct {
		name string
		file string
		cmds []string
	}{
		{
			name: "reset",
			file: "./test-fixtures/reset.json",
			cmds: []string{},
		},
	}
	for _, tt := range tests {
		ctx := newApp(tt.file, tt.cmds)
		lbSetup(tt.file)
		Reset(ctx)
		lb := Read(ctx)
		for _, entry := range lb.Entries {
			if entry.Points != 0 {
				log.Fatalf("%q. entry %v = %v, want 0", tt.name, entry.Name, entry.Points)
			}
		}
	}
}

func TestRead(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{
			name: "read",
			file: "./test-fixtures/read.json",
		},
	}
	for _, tt := range tests {
		want := lbSetup(tt.file)

		if got := Read(newApp(tt.file, []string{})); !reflect.DeepEqual(got, want) {
			t.Errorf("%q. Read() = %v, want %v", tt.name, got, want)
		}
	}
}

func currPoints(lb *Leaderboard, name string) int {
	for _, entry := range lb.Entries {
		if strings.EqualFold(name, entry.Name) {
			return entry.Points
		}
	}
	return 0
}

func newApp(file string, cmds []string) *cli.Context {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Value: file,
		},
	}
	set := flag.NewFlagSet("test", 3)
	set.String("file", file, "")
	if err := set.Parse(cmds); err != nil {
		log.Fatalf("error on setup, %s", err.Error())
	}
	return cli.NewContext(app, set, nil)
}

func lbSetup(key string) *Leaderboard {
	//detect if file exists
	var _, err = os.Stat(key)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(key)
		if err != nil {
			log.Fatalf("error on setup of test leaderboard, key: %s, error: %s", key, err.Error())
		}
		_ = file.Close()
	}
	lb := &Leaderboard{
		Key:     key,
		Headers: []string{"name", "points"},
		Entries: []*Entry{
			{
				Name:   "Gordon",
				Points: 3,
			},
			{
				Name:   "Jane",
				Points: 2,
			},
			{
				Name:   "Slackbot",
				Points: 1,
			},
		},
	}
	lb.Save()
	return lb
}
