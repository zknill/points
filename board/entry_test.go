package board

import (
	"testing"

	"github.com/pkg/errors"
	"google.golang.org/appengine/datastore"
)

func Test_aeEntry_Score(t *testing.T) {
	tests := []struct {
		name    string
		aee     *aeEntry
		want    int
		wantErr bool
	}{
		{name: "alice - 3", aee: &aeEntry{entryKey("alice")}, want: 3, wantErr: false},
		{name: "bob - 2", aee: &aeEntry{entryKey("bob")}, want: 2, wantErr: false},
		{name: "jane - 1", aee: &aeEntry{entryKey("jane")}, want: 1, wantErr: false},
		{name: "missing entry", aee: &aeEntry{entryKey("slackbot")}, want: 0, wantErr: true},
		{name: "error key", aee: &aeEntry{datastore.NewKey(aectx, "Unknown", "error", 0, nil)}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.aee.Score(aectx)
			if (err != nil) != tt.wantErr {
				t.Errorf("aeEntry.Score() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("aeEntry.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aeEntry_Add(t *testing.T) {
	tests := []struct {
		name    string
		aee     *aeEntry
		num     int
		before  int
		wantErr bool
	}{
		{name: "add 1 alice", aee: &aeEntry{entryKey("alice")}, num: 1, before: 3},
		{name: "add 4 bob", aee: &aeEntry{entryKey("bob")}, num: 4, before: 2},
		{name: "add 0 jane", aee: &aeEntry{entryKey("jane")}, num: 0, before: 1},
		{name: "add 1 slackbot", aee: &aeEntry{entryKey("slackbot")}, num: 1, before: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBefore, err := tt.aee.Score(aectx)
			if err != nil && errors.Cause(err) != ErrEntryNotFound {
				t.Fatalf("%q. got error check score before, %+v", tt.name, err)
			}
			if gotBefore != tt.before {
				t.Errorf("%q. incorrect score before, want = %v, got = %v", tt.name, tt.before, gotBefore)
			}
			if err := tt.aee.Add(aectx, tt.num); err != nil {
				t.Fatalf("%q. error adding points, %v", tt.name, err)
			}
			gotAfter, err := tt.aee.Score(aectx)
			if err != nil {
				t.Fatalf("%q. got error check score after, %+v", tt.name, err)
			}
			if gotBefore != tt.before {
				t.Errorf("%q. incorrect score after, want = %v, got = %v", tt.name, tt.before+tt.num, gotAfter)
			}
		})
	}
}
