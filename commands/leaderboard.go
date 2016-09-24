package points

import (
	"io/ioutil"
	"encoding/csv"
	"strings"
	"io"
	"strconv"
	"os"
)

type Leaderboard struct {
	Headers  []string
	Entries  []*Entry
	History  []*History
	filename string
}

type History struct {
	Timestamp int
	Message   string
}

func (lb *Leaderboard) Load(filename string) {
	lb.filename = filename
	in, _ := ioutil.ReadFile(filename)
	r := csv.NewReader(strings.NewReader(string(in)))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		checkErr(err)
		points, err := strconv.Atoi(record[1])
		if err != nil && lb.Headers == nil {
			lb.Headers = []string{record[0], record[1]}
			continue
		}
		lb.Entries = append(lb.Entries, &Entry{record[0], points})
	}
}

func (lb *Leaderboard) Save() {
	file, err := os.Create(filename)
	checkErr(err)
	defer file.Close()

	writer := csv.NewWriter(file)

	err = writer.Write(lb.Headers)
	checkErr(err)
	for _, entry := range lb.Entries {
		err := writer.Write(entry.Array())
		checkErr(err)
	}

	defer writer.Flush()
}
