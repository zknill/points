package points

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Leaderboard struct {
	Headers []string
	Entries []*Entry
	History []*History
	Key     string
}

type StoredLeaderboard struct {
	Headers []string
}

func (lb *Leaderboard) Load(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatal("Leaderboard '" + filename + "' not yet initialised!")
	}

	lb.Key = filename
	in, _ := ioutil.ReadFile(filename)
	json.Unmarshal(in, &lb)
}

func (lb *Leaderboard) Save() {
	b, _ := json.Marshal(lb)
	file, err := os.Create(lb.Key)
	checkErr(err)
	defer file.Close()

	file.Write(b)
}

func (lb *Leaderboard) Add(name, pnts string, meta []string) {
	var err error
	var number = 0
	if pnts != "" {
		number, err = strconv.Atoi(pnts)
		checkErr(err)
	}
	found := false
	for _, entry := range lb.Entries {
		if strings.EqualFold(name, entry.Name) {
			found = true
			entry.Points += number
			if len(lb.Headers) > 2 {
				newMeta := meta[:len(lb.Headers)-2]
				if newMeta == nil {
					newMeta = []string{}
				}
				entry.Meta = newMeta
			}
		}
	}
	if !found {
		lb.Entries = append(lb.Entries, &Entry{strings.Title(name), number, meta})
	}
}
