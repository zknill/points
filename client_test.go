package points

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestClient_Entry(t *testing.T) {
	tests := []struct {
		name          string
		team          string
		entryName     string
		expectedEntry *Entry
		entryError    error
	}{
		{name: "success", expectedEntry: &Entry{Name: "slackbot", Score: 0}, entryName: "slackbot", team: "slackteam"},
		{name: "failure", entryError: errors.New("entry not found"), entryName: "jirabot", team: "atlassian"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := nopStore{
				entry: func(name TeamName, entryName EntryName) (*Entry, error) {
					if name != TeamName(tt.team) {
						t.Errorf("nopStore.Entry(%s, %s) wanted team name %s", name, entryName, tt.team)
					}

					if entryName != EntryName(tt.entryName) {
						t.Errorf("notStore.Entry(%s, %s) wanted entry name %s", name, entryName, tt.entryName)
					}

					return tt.expectedEntry, tt.entryError
				},
			}

			client := &Client{
				team:      TeamName(tt.team),
				store:     store,
				log:       defaultLogger(),
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

func TestClient_Scores(t *testing.T) {
	type fields struct {
		team      TeamName
		store     entryStore
		formatter Formatter
		log       logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantW   string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				team:      tt.fields.team,
				store:     tt.fields.store,
				formatter: tt.fields.formatter,
				log:       tt.fields.log,
			}
			w := &bytes.Buffer{}
			if err := c.Scores(tt.args.ctx, w); (err != nil) != tt.wantErr {
				t.Errorf("Client.Scores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Client.Scores() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestClient_NewEntry(t *testing.T) {
	type fields struct {
		team      TeamName
		store     entryStore
		formatter Formatter
		log       logger
	}
	type args struct {
		ctx   context.Context
		entry *Entry
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				team:      tt.fields.team,
				store:     tt.fields.store,
				formatter: tt.fields.formatter,
				log:       tt.fields.log,
			}
			if err := c.NewEntry(tt.args.ctx, tt.args.entry); (err != nil) != tt.wantErr {
				t.Errorf("Client.NewEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Add(t *testing.T) {
	type fields struct {
		team      TeamName
		store     entryStore
		formatter Formatter
		log       logger
	}
	type args struct {
		ctx  context.Context
		name EntryName
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				team:      tt.fields.team,
				store:     tt.fields.store,
				formatter: tt.fields.formatter,
				log:       tt.fields.log,
			}
			if err := c.Add(tt.args.ctx, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Client.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Reset(t *testing.T) {
	type fields struct {
		team      TeamName
		store     entryStore
		formatter Formatter
		log       logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				team:      tt.fields.team,
				store:     tt.fields.store,
				formatter: tt.fields.formatter,
				log:       tt.fields.log,
			}
			if err := c.Reset(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Client.Reset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_getEntry(t *testing.T) {
	type fields struct {
		team      TeamName
		store     entryStore
		formatter Formatter
		log       logger
	}
	type args struct {
		ctx  context.Context
		name EntryName
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Entry
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				team:      tt.fields.team,
				store:     tt.fields.store,
				formatter: tt.fields.formatter,
				log:       tt.fields.log,
			}
			got, err := c.getEntry(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.getEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.getEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_putEntry(t *testing.T) {
	type fields struct {
		team      TeamName
		store     entryStore
		formatter Formatter
		log       logger
	}
	type args struct {
		ctx   context.Context
		entry *Entry
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				team:      tt.fields.team,
				store:     tt.fields.store,
				formatter: tt.fields.formatter,
				log:       tt.fields.log,
			}
			if err := c.putEntry(tt.args.ctx, tt.args.entry); (err != nil) != tt.wantErr {
				t.Errorf("Client.putEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_getAll(t *testing.T) {
	type fields struct {
		team      TeamName
		store     entryStore
		formatter Formatter
		log       logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Entry
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				team:      tt.fields.team,
				store:     tt.fields.store,
				formatter: tt.fields.formatter,
				log:       tt.fields.log,
			}
			got, err := c.getAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.getAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.getAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_runTransaction(t *testing.T) {
	type fields struct {
		team      TeamName
		store     entryStore
		formatter Formatter
		log       logger
	}
	type args struct {
		ctx context.Context
		tx  transaction
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				team:      tt.fields.team,
				store:     tt.fields.store,
				formatter: tt.fields.formatter,
				log:       tt.fields.log,
			}
			if err := c.runTransaction(tt.args.ctx, tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("Client.runTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ensureTeamExists(t *testing.T) {
	type fields struct {
		team      TeamName
		store     entryStore
		formatter Formatter
		log       logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				team:      tt.fields.team,
				store:     tt.fields.store,
				formatter: tt.fields.formatter,
				log:       tt.fields.log,
			}
			if err := c.ensureTeamExists(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Client.ensureTeamExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientFactory_New(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		cf      *ClientFactory
		args    args
		want    *Client
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf := &ClientFactory{}
			got, err := cf.New(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientFactory.New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientFactory.New() = %v, want %v", got, tt.want)
			}
		})
	}
}

type nopStore struct {
	entry func(name TeamName, entryName EntryName) (*Entry, error)
}

func (n nopStore) Entry(ctx context.Context, team TeamName, name EntryName) (*Entry, error) {
	return n.entry(team, name)
}

func (nopStore) Entries(ctx context.Context, team TeamName) ([]*Entry, error) {
	panic("implement me")
}

func (nopStore) PutEntry(ctx context.Context, team TeamName, entry *Entry) error {
	panic("implement me")
}

func (nopStore) Team(ctx context.Context, team TeamName) (*Team, error) {
	panic("implement me")
}

func (nopStore) PutTeam(ctx context.Context, team *Team) error {
	panic("implement me")
}
