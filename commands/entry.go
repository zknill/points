package points

import "strconv"

// Entry a single entry in the table
type Entry struct {
	Name   string
	Points int
}

// ScoreFirst sorts by score then name
type ScoreFirst []*Entry

// Len implements sort interface
func (e ScoreFirst) Len() int {
	return len(e)
}

// Swap implements sort interface
func (e ScoreFirst) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// Less implements sort interface
func (e ScoreFirst) Less(i, j int) bool {
	if e[j].Points > 0 || e[i].Points > 0 {
		return e[i].Points > e[j].Points
	}
	return e[i].Name < e[j].Name
}

// String formats the entry
func (e *Entry) String() string {
	score := strconv.Itoa(e.Points)
	return e.Name + ": " + score
}

// Array entry as a string array
func (e *Entry) Array() []string {
	return []string{e.Name, strconv.Itoa(e.Points)}
}
