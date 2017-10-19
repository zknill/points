package points

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/zknill/points/internal/log"
	"github.com/zknill/points/internal/store"
)

func TestClient_Entry(t *testing.T) {
	tests := []struct {
		name          string
		team          string
		entryName     string
		storeEntry    *store.Entry
		expectedEntry *Entry
		entryError    error
	}{
		{name: "success",
			storeEntry:    &store.Entry{Name: "slackbot", Score: 0},
			expectedEntry: &Entry{Name: "slackbot", Score: 0},
			entryName:     "slackbot",
			team:          "slackteam"},
		{name: "failure", entryError: errors.New("entry not found"), entryName: "jirabot", team: "atlassian"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := nopStore{
				entry: func(name string, entryName string) (*store.Entry, error) {
					if name != tt.team {
						t.Errorf("nopStore.Entry(%s, %s) wanted team name %s", name, entryName, tt.team)
					}

					if entryName != tt.entryName {
						t.Errorf("notStore.Entry(%s, %s) wanted entry name %s", name, entryName, tt.entryName)
					}

					return tt.storeEntry, tt.entryError
				},
			}

			client := &TeamClient{
				team:      TeamName(tt.team),
				store:     ns,
				log:       log.DefaultLogger(),
				formatter: formatFunc(slackFormatter),
			}

			ctx := context.Background()
			got, err := client.Entry(ctx, EntryName(tt.entryName))
			if !reflect.DeepEqual(err, tt.entryError) {
				t.Errorf("client.Entry(ctx, %s) got error does not match expected, got: %+v, expected: %+v", tt.entryName, err, tt.entryError)
			}

			if !reflect.DeepEqual(got, tt.expectedEntry) {
				t.Errorf("client.Entry(ctx, %s) got entry does match expected, got: %+v, expected: %+v", tt.entryName, got, tt.expectedEntry)
			}
		})
	}
}

var _ entryStore = (*nopStore)(nil)

type nopStore struct {
	entry func(name string, entryName string) (*store.Entry, error)
}

func (n nopStore) Entry(ctx context.Context, team string, name string) (*store.Entry, error) {
	return n.entry(team, name)
}

func (nopStore) Entries(ctx context.Context, team string) ([]*store.Entry, error) {
	panic("implement me")
}

func (nopStore) PutEntry(ctx context.Context, team string, entry *store.Entry) error {
	panic("implement me")
}

func (nopStore) Team(ctx context.Context, team string) (*store.Team, error) {
	panic("implement me")
}

func (nopStore) PutTeam(ctx context.Context, team *store.Team) error {
	panic("implement me")
}
