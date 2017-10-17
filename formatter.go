package points

import (
	"io"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

type Formatter interface {
	FormatEntries(w io.Writer, headers []string, entries []*Entry)
}

type formatFunc func(w io.Writer, headers []string, entries []*Entry)

func (ff formatFunc) FormatEntries(w io.Writer, headers []string, entries []*Entry) {
	ff(w, headers, entries)
}

func slackFormatter(w io.Writer, headers []string, entries []*Entry) {
	sort.Slice(entries, func(i, j int) bool {
		iEntry, jEntry := entries[i], entries[j]
		if iEntry.Score == jEntry.Score {
			return iEntry.Name < jEntry.Name
		}
		return iEntry.Score > jEntry.Score
	})

	tw := tablewriter.NewWriter(w)
	tw.SetHeader(headers)

	for _, entry := range entries {
		tw.Append([]string{entry.Name, strconv.Itoa(entry.Score)})
	}

	tw.Render()
}
