package points

import (
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"fmt"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

const filename = "points.json"

func Print(_ *cli.Context) {
	lb := read()
	printTable(lb.Headers, lb.Entries)
}

func Add(c *cli.Context) {
	lb := read()
	name := c.Args().Get(0)
	if name == "" {
		return
	}

	arg1 := c.Args().Get(1)
	var err error
	var number = 0
	if arg1 != "" {
		number, err = strconv.Atoi(arg1)
		checkErr(err)
	}
	found := false
	for _, entry := range lb.Entries {
		if strings.ToUpper(name) == strings.ToUpper(entry.Name) {
			found = true
			entry.Points += number
		}
	}
	if !found {
		lb.Entries = append(lb.Entries, &Entry{strings.Title(name), number})
	}
	lb.Save()
}

func Reset(_ *cli.Context) {
	lb := read()
	for _, entry := range lb.Entries {
		entry.Points = 0
	}
	lb.Save()
}

func Slack(_ *cli.Context) {
	fmt.Println("```")
	lb := read()
	printTable(lb.Headers, lb.Entries)
	fmt.Println("```")
}

func Init(_ *cli.Context) {
	lb := &Leaderboard{}
	lb.filename = filename
	lb.Save()
}

func read() *Leaderboard {
	lb := &Leaderboard{}
	lb.Load(filename)
	return lb
}

func printTable(headers []string, entries []*Entry) {
	sort.Sort(PointsFirst(entries))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	for _, entry := range entries {
		table.Append(entry.Array())
	}
	table.Render()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
