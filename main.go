package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/zknill/points/commands"
)

func main() {
	app := cli.NewApp()
	app.Name = "points"
	app.Usage = "the leaderboard!"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "show the results",
			UsageText: "shows the leader board printed in an ASCII table " +
				"and ordered by points and then name",
			Action: points.Print,
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "adds points to a team member",
			UsageText: "adds a team member or points. " +
				"team memebers can be added with zero points. " +
				"points can be added to a team members existing total",
			Action: points.Add,
		},
		{
			Name:    "reset",
			Aliases: []string{"r"},
			Usage:   "resets all the points of the team members or specific member back to 0",
			Action:  points.Reset,
			Flags: []cli.Flag {
				cli.StringFlag{
					Name: "entry",
					Value: "all",
					Usage: "entry to reset",
				},
			},
		},
		{
			Name:    "slack",
			Aliases: []string{"s"},
			Usage:   "prints point totals in a slack friendly way",
			Action:  points.Slack,
		},
		{
			Name:      "init",
			Aliases:   []string{"i"},
			Usage:     "creates a new json file backend to store the leaderboard",
			UsageText: "points init headers...",
			Action:    points.Init,
		},
		{
			Name:    "history",
			Aliases: []string{"p"},
			Usage:   "shows the command history",
			Action:  points.ShowHistory,
		},
	}
	app.Run(os.Args)
}
