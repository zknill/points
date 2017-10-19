package store

import (
	"context"

	"google.golang.org/appengine/datastore"
)

const (
	kindTeam  string = "team"
	kindEntry string = "entry"
)

// AppEngine defines methods for interacting with the google app
// engine datastore. It has all the methods needed by points.
type AppEngine struct{}

// Entry represents the fields that are
// stored in the datastore for an entry.
type Entry struct {
	Name  string
	Score int
}

// Team is the fields that are stored
// in the datastore for a team.
type Team struct {
	Name    string
	Headers []string
}

// Transaction defines a transaction. Transactions are passed a transaction
// context and the transaction is considered failed if the returned error
// is non-nil.
type Transaction func(tc context.Context) error

// Transact defines a method that can run a transaction against the
// google app engine datastore.
func (ae *AppEngine) Transact(ctx context.Context, tx Transaction) error {
	return datastore.RunInTransaction(ctx, tx, nil)
}

// Entry returns a single points.Entry for a team and entry name.
func (ae *AppEngine) Entry(ctx context.Context, team, name string) (*Entry, error) {
	entryKey := ae.entryKey(ctx, team, name)

	var entry = new(Entry)

	return entry, datastore.Get(ctx, entryKey, entry)
}

// Entries returns all entries for a team
func (ae *AppEngine) Entries(ctx context.Context, team string) ([]*Entry, error) {
	teamKey := ae.teamKey(ctx, team)

	query := datastore.NewQuery(kindEntry).Ancestor(teamKey)
	results := query.Run(ctx)

	var entries []*Entry
	for {
		entry := new(Entry)

		if _, err := results.Next(entry); err != nil {
			if err == datastore.Done {
				return entries, nil
			}
			return nil, err
		}

		entries = append(entries, entry)
	}
}

// PutEntry puts a single entry and associates it as a child of the team.
func (ae *AppEngine) PutEntry(ctx context.Context, team string, entry *Entry) error {
	entryKey := ae.entryKey(ctx, team, entry.Name)

	_, err := datastore.Put(ctx, entryKey, entry)

	return err
}

// Team get the team by name.
func (ae *AppEngine) Team(ctx context.Context, teamName string) (*Team, error) {
	teamKey := ae.teamKey(ctx, teamName)

	var team = new(Team)
	err := datastore.Get(ctx, teamKey, team)

	return team, err
}

// PutTeam stores a team by it's team name.
func (ae *AppEngine) PutTeam(ctx context.Context, team *Team) error {
	teamKey := ae.teamKey(ctx, team.Name)

	_, err := datastore.Put(ctx, teamKey, team)
	return err
}

func (ae *AppEngine) teamKey(ctx context.Context, team string) *datastore.Key {
	return datastore.NewKey(ctx, kindTeam, team, 0, nil)
}

func (ae *AppEngine) entryKey(ctx context.Context, team, name string) *datastore.Key {
	teamKey := ae.teamKey(ctx, team)
	return datastore.NewKey(ctx, kindEntry, name, 0, teamKey)
}
