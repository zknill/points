package points

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

// Filename represents where to store the cli app's backend file
const Filename = "./points.json"

// Print leaderboard to console
func Print(c *cli.Context) {
	lb := Read(c)
	PrintTable(nil, lb)
}

// Add points to members
func Add(c *cli.Context) {
	lb := Read(c)
	name := c.Args().Get(0)
	if name == "" {
		return
	}
	var pnts = "1"
	if len(c.Args()) > 1 {
		arg1 := c.Args().Get(1)
		if arg1 != "" {
			pnts = arg1
		}
	}
	_ = lb.Add(name, pnts)
	lb.addHistory("add", name, pnts)
	lb.Save()
}

// Reset points for all or a single member
func Reset(c *cli.Context) {
	lb := Read(c)
	flag := c.String("entry")
	if flag == "all" || flag == "" {
		for _, entry := range lb.Entries {
			entry.Points = 0
		}
	} else {
		for _, entry := range lb.Entries {
			if strings.EqualFold(entry.Name, flag) {
				entry.Points = 0
				break
			}
		}
	}
	lb.addHistory("reset")
	lb.Save()
}

// Slack print table in a slack friendly way
func Slack(c *cli.Context) {
	fmt.Println("```")
	lb := Read(c)
	PrintTable(nil, lb)
	fmt.Println("```")
}

// InitStorage creates a new storage file backend
func InitStorage(c *cli.Context) {
	lb := &Leaderboard{}
	lb.Key = Filename
	if _, err := os.Stat(lb.Key); os.IsNotExist(err) {
		var headers []string
		if headers = c.Args(); headers[0] == "" {
			headers = []string{"name", "points"}
		}
		lb.Headers = headers
		lb.addHistory("init", headers...)
		lb.Save()
		return
	}
	log.Fatal("A leaderboard already exists in this directory")
}

// ShowHistory prints the history to console
func ShowHistory(c *cli.Context) {
	lb := Read(c)
	for _, h := range lb.History {
		fmt.Println(h.String())
	}
}

// Read loads the leaderboard from file into memory
func Read(c *cli.Context) *Leaderboard {
	lb := &Leaderboard{Key: Filename}
	if file := c.String("file"); file != "" {
		lb.Key = file
	}
	lb.Load()
	return lb
}

// PrintTable prints the table to console
func PrintTable(_ context.Context, lb *Leaderboard) {
	table := GetTable(os.Stdout, lb)
	table.Render()
}

// GetTable the table from the leaderboard
func GetTable(writer io.Writer, lb *Leaderboard) *tablewriter.Table {
	sort.Sort(ScoreFirst(lb.Entries))
	table := tablewriter.NewWriter(writer)
	table.SetHeader(lb.Headers)
	for _, entry := range lb.Entries {
		table.Append(entry.Array())
	}
	return table
}
