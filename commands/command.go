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

	var number = 0

	if len(c.Args()) == 1 {
		number = 1
	} else {
		arg1 := c.Args().Get(1)
		var err error
		
		if arg1 != "" {
			number, err = strconv.Atoi(arg1)
			checkErr(err)
		}
	}

	found := false
	for _, entry := range lb.Entries {
		if strings.EqualFold(name, entry.Name) {
			found = true
			entry.Points += number
			if len(lb.Headers)> 2 {
				meta := meta(c)[:len(lb.Headers) - 2]
				if meta == nil {
					meta = []string{}
				}
				entry.Meta = meta
			}
		}
	}
	if !found {
		lb.Entries = append(lb.Entries, &Entry{strings.Title(name), number, meta(c)[:len(lb.Headers) - 2]})
	}
	lb.addHistory("add", name, strconv.Itoa(number))
	lb.Save()
}

func Reset(c *cli.Context) {
	lb := read()
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
	lb := read()
	printTable(lb)
	fmt.Println("```")
}

func Init(c *cli.Context) {
	lb := &Leaderboard{}
	lb.filename = filename
	if _, err := os.Stat(lb.filename); os.IsNotExist(err) {
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

func args(c *cli.Context) []string {
	return append([]string{c.Args().First()}, c.Args().Tail()...)
}

func meta(c *cli.Context) []string {
	meta := []string{}
	if len(c.Args().Tail()) > 1 {
		meta = c.Args().Tail()[1:]
	}
	return meta
}
