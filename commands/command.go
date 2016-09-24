package points

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/olekukonko/tablewriter"
	"fmt"
)

const filename = "points.csv"

func Print(_ *cli.Context) {
	printTable(read())
}

func Add(c *cli.Context) {
	headers, entries := read()
	name := c.Args().Get(0)

	arg1 := c.Args().Get(1)
	var number int
	var err error
	if arg1 != "" {
		number, err = strconv.Atoi(arg1)
		checkErr(err)
	} else {
		number = 0
	}
	found := false
	for _, entry := range entries {
		if strings.ToUpper(name) == strings.ToUpper(entry.Name) {
			found = true
			entry.Points += number
		}
	}
	if !found {
		entries = append(entries, &Entry{strings.Title(name), number})
	}
	printTable(headers, entries)
	saveTable(headers, entries)
}

func Reset(_ *cli.Context) {
	headers, entries := read()
	for _, entry := range entries {
		entry.Points = 0
	}
	printTable(headers, entries)
	saveTable(headers, entries)
}

func Slack(_ *cli.Context) {
	fmt.Println("```")
	printTable(read())
	fmt.Println("```")
}

func read() ([]string, []*Entry) {
	in, _ := ioutil.ReadFile(filename)
	r := csv.NewReader(strings.NewReader(string(in)))
	entries := []*Entry{}
	var headers []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		checkErr(err)
		points, err := strconv.Atoi(record[1])
		if err != nil {
			headers = []string{record[0], record[1]}
			continue
		}
		entries = append(entries, &Entry{record[0], points})
	}
	return headers, entries
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

func saveTable(headers []string, entries []*Entry) {
	file, err := os.Create(filename)
	checkErr(err)
	defer file.Close()

	writer := csv.NewWriter(file)

	err = writer.Write(headers)
	checkErr(err)
	for _, entry := range entries {
		err := writer.Write(entry.Array())
		checkErr(err)
	}

	defer writer.Flush()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
