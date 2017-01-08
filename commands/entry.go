package points

import "strconv"

type Entry struct {
	Name   string
	Points int
}

type ScoreFirst []*Entry

func (e ScoreFirst) Len() int {
	return len(e)
}

func (e ScoreFirst) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e ScoreFirst) Less(i, j int) bool {
	if e[j].Points > 0 || e[i].Points > 0 {
		return e[i].Points > e[j].Points
	}
	return e[i].Name < e[j].Name
}

func (e *Entry) String() string {
	score := strconv.Itoa(e.Points)
	return e.Name + ": " + score
}

func (e *Entry) Array() []string {
	return append([]string{e.Name, strconv.Itoa(e.Points)})
}
