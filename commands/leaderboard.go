package points

import (
	"fmt"
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

func (lb *Leaderboard) Load(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "A leaderboard '%v' does not exist in this directory!\n", filename)
	}

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
