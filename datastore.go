package points

import (
	"context"

	"google.golang.org/appengine/datastore"
)

const (
	kindTeam  string = "team"
	kindEntry string = "entry"
)

type entryStore interface {
	Entry(ctx context.Context, team TeamName, name EntryName) (*Entry, error)
	Entries(ctx context.Context, team TeamName) ([]*Entry, error)
	PutEntry(ctx context.Context, team TeamName, entry *Entry) error

	Team(ctx context.Context, team TeamName) (*Team, error)
	PutTeam(ctx context.Context, team *Team) error
}

var _ entryStore = (*appEngineDatastore)(nil)

type transaction func(tc context.Context) error

type transacter interface {
	transact(tx transaction) error
}

type appEngineDatastore struct{}

func (ae *appEngineDatastore) transact(ctx context.Context, tx transaction) error {
	return datastore.RunInTransaction(ctx, tx, nil)
}

func (ae *appEngineDatastore) Entry(ctx context.Context, team TeamName, name EntryName) (*Entry, error) {
	entryKey := ae.entryKey(ctx, team, name)

	var entry = new(Entry)
	tx := func(tc context.Context) error {
		return datastore.Get(ctx, entryKey, entry)
	}

	err := datastore.RunInTransaction(ctx, tx, nil)
	return entry, err
}

func (ae *appEngineDatastore) Entries(ctx context.Context, team TeamName) ([]*Entry, error) {
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

func (ae *appEngineDatastore) PutEntry(ctx context.Context, team TeamName, entry *Entry) error {
	entryKey := ae.entryKey(ctx, team, EntryName(entry.Name))

	tx := func(tc context.Context) error {
		_, err := datastore.Put(ctx, entryKey, entry)
		return err
	}

	return datastore.RunInTransaction(ctx, tx, nil)
}

func (ae *appEngineDatastore) Team(ctx context.Context, teamName TeamName) (*Team, error) {
	teamKey := ae.teamKey(ctx, teamName)

	var team = new(Team)
	err := datastore.Get(ctx, teamKey, team)

	return team, err
}

func (ae *appEngineDatastore) PutTeam(ctx context.Context, team *Team) error {
	teamKey := ae.teamKey(ctx, TeamName(team.Name))

	_, err := datastore.Put(ctx, teamKey, team)
	return err
}

func (ae *appEngineDatastore) teamKey(ctx context.Context, team TeamName) *datastore.Key {
	return datastore.NewKey(ctx, kindTeam, team.String(), 0, nil)
}

func (ae *appEngineDatastore) entryKey(ctx context.Context, team TeamName, name EntryName) *datastore.Key {
	teamKey := ae.teamKey(ctx, team)
	return datastore.NewKey(ctx, kindEntry, name.String(), 0, teamKey)
}
