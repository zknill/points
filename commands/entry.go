package points

import "strconv"

type Entry struct {
	Name   string
	Points int
	Meta []string
}

type PointsFirst []*Entry

func (e PointsFirst) Len() int {
	return len(e)
}

func (e PointsFirst) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e PointsFirst) Less(i, j int) bool {
	if e[j].Points > 0 || e[i].Points > 0 {
		return e[i].Points > e[j].Points
	}
	return e[i].Name < e[j].Name
}

func (e *Entry) String() string {
	pointsStr := strconv.Itoa(e.Points)
	return e.Name + ": " + pointsStr
}

func (e *Entry) Array() []string {
	return append([]string{e.Name, strconv.Itoa(e.Points)}, e.Meta...)
}
