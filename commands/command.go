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

const Filename = "/Users/zak/gopath/src/github.com/zknill/points/points.json"

func Print(_ *cli.Context) {
	lb := Read()
	PrintTable(nil, lb)
}

func Add(c *cli.Context) {
	lb := Read()
	name := c.Args().Get(0)
	if name == "" {
		return
	}
	pnts := c.Args().Get(1)
	lb.Add(name, pnts)
	lb.addHistory("add", args(c)[:len(lb.Headers)]...)
	lb.Save()
}

func Reset(c *cli.Context) {
	lb := Read()
	flag := c.String("entry")
	if flag == "all" {
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

func Slack(_ *cli.Context) {
	fmt.Println("```")
	lb := Read()
	PrintTable(nil, lb)
	fmt.Println("```")
}

func InitStorage(c *cli.Context) {
	lb := &Leaderboard{}
	lb.Key = Filename
	if _, err := os.Stat(lb.Key); os.IsNotExist(err) {
		var headers []string
		if headers = args(c); headers[0] == "" {
			headers = []string{"name", "points"}
		}
		lb.Headers = headers
		lb.addHistory("init", headers...)
		lb.Save()
		return
	}
	log.Fatal("A leaderboard already exists in this directory")
}

func ShowHistory(c *cli.Context) {
	lb := Read()
	for _, h := range lb.History {
		fmt.Println(h.String())
	}
}

func Read() *Leaderboard {
	lb := &Leaderboard{}
	lb.Load(Filename)
	return lb
}

func PrintTable(_ context.Context, lb *Leaderboard) {
	table := GetTable(os.Stdout, lb)
	table.Render()
}

func GetTable(writer io.Writer, lb *Leaderboard) *tablewriter.Table {
	sort.Sort(ScoreFirst(lb.Entries))
	table := tablewriter.NewWriter(writer)
	table.SetHeader(lb.Headers)
	for _, entry := range lb.Entries {
		table.Append(entry.Array())
	}
	return table
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func args(c *cli.Context) []string {
	return append([]string{c.Args().First()}, c.Args().Tail()...)
}
