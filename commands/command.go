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
	printTable(lb)
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
	lb.addHistory("add", getArgs(c)...)
	lb.Save()
}

func Reset(_ *cli.Context) {
	lb := read()
	for _, entry := range lb.Entries {
		entry.Points = 0
	}
	lb.addHistory("reset")
	lb.Save()
}

func Slack(_ *cli.Context) {
	fmt.Println("```")
	lb := read()
	printTable(lb)
	fmt.Println("```")
}

func Init(c *cli.Context) {
	lb := &Leaderboard{}
	lb.filename = filename
	if _, err := os.Stat(lb.filename); os.IsNotExist(err) {
		lb.Headers = append([]string{c.Args().First()}, c.Args().Tail()...)
		lb.addHistory("init", getArgs(c)...)
		lb.Save()
		return
	}
	log.Fatal("A leaderboard already exists in this directory")
}

func ShowHistory(c *cli.Context) {
	lb := read()
	for _, h := range lb.History {
		fmt.Println(h.string())
	}
}

func read() *Leaderboard {
	lb := &Leaderboard{}
	lb.Load(filename)
	return lb
}

func printTable(lb *Leaderboard) {
	sort.Sort(PointsFirst(lb.Entries))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(lb.Headers)
	for _, entry := range lb.Entries {
		table.Append(entry.Array())
	}
	table.Render()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getArgs(c *cli.Context) []string {
	return append([]string{c.Args().First()}, c.Args().Tail()...)
}
