package points

import (
	"encoding/json"
	"io/ioutil"
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
	json.Unmarshal(in, &lb)
}

func (lb *Leaderboard) Save() {
	b, _ := json.Marshal(lb)
	file, err := os.Create(lb.filename)
	checkErr(err)
	defer file.Close()

	file.Write(b)
}
