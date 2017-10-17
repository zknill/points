package points

import (
	"context"
	"io"

	"google.golang.org/appengine/datastore"
)

var defaultHeaders = []string{"name", "points"}

type EntryName string

func (n EntryName) String() string {
	return string(n)
}

type Entry struct {
	Name  string
	Score int
}

type TeamName string

func (t TeamName) String() string {
	return string(t)
}

type Team struct {
	Name    string
	Headers []string
}

type ClientAdapter func(client *Client) *Client

type Client struct {
	team      TeamName
	store     entryStore
	formatter Formatter
	log       logger
}

func New(team TeamName, adapters ...ClientAdapter) *Client {
	client := &Client{
		team: team,
		log:  defaultLogger(),
	}

	for _, adapter := range adapters {
		client = adapter(client)
	}

	return client
}

func WithAppEngine() ClientAdapter {
	return func(client *Client) *Client {
		client.log = &appEngineLog{}
		client.store = &appEngineDatastore{}

		return client
	}
}

func WithSlackFormatter() ClientAdapter {
	return func(client *Client) *Client {
		client.formatter = formatFunc(slackFormatter)
		return client
	}
}

func (c *Client) Entry(ctx context.Context, name EntryName) (*Entry, error) {
	var entry = new(Entry)
	var tx transaction = func(tc context.Context) error {
		var err error
		entry, err = c.getEntry(tc, name)
		return err
	}

	return entry, c.runTransaction(ctx, tx)
}

func (c *Client) Scores(ctx context.Context, w io.Writer) error {
	var err error
	var team *Team
	var entries []*Entry

	var tx transaction = func(tc context.Context) error {
		team, err = c.store.Team(ctx, c.team)
		if err != nil {
			return err
		}

		entries, err = c.getAll(ctx)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.runTransaction(ctx, tx)
	if err != nil {
		return err
	}

	c.formatter.FormatEntries(w, team.Headers, entries)
	return nil
}

func (c *Client) NewEntry(ctx context.Context, entry *Entry) error {
	var tx transaction = func(tc context.Context) error {
		return c.putEntry(tc, entry)
	}

	return c.runTransaction(ctx, tx)
}

func (c *Client) Add(ctx context.Context, name EntryName) error {

	var tx transaction = func(tc context.Context) error {

		entry, err := c.getEntry(ctx, name)
		if err != nil {

			if err != datastore.ErrNoSuchEntity {
				return err
			}

			entry = &Entry{Name: name.String(), Score: 1}
		}

		entry.Score += 1

		return c.putEntry(ctx, entry)
	}

	return c.runTransaction(ctx, tx)
}

func (c *Client) Reset(ctx context.Context) error {

	var tx transaction = func(tc context.Context) error {

		entries, err := c.getAll(tc)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			entry.Score = 0
			if err := c.putEntry(ctx, entry); err != nil {
				return err
			}
		}

		return nil
	}

	return c.runTransaction(ctx, tx)
}

func (c *Client) getEntry(ctx context.Context, name EntryName) (*Entry, error) {
	return c.store.Entry(ctx, c.team, name)
}

func (c *Client) putEntry(ctx context.Context, entry *Entry) error {
	return c.store.PutEntry(ctx, c.team, entry)
}

func (c *Client) getAll(ctx context.Context) ([]*Entry, error) {
	return c.store.Entries(ctx, c.team)
}

func (c *Client) runTransaction(ctx context.Context, tx transaction) error {
	if d, ok := c.store.(transacter); ok {
		return d.transact(tx)
	}
	return tx(ctx)
}

func (c *Client) ensureTeamExists(ctx context.Context) error {
	var tx transaction = func(tc context.Context) error {
		_, err := c.store.Team(ctx, c.team)
		if err == nil {
			return nil
		}

		if err != datastore.ErrNoSuchEntity {
			return err
		}

		return c.store.PutTeam(ctx, &Team{Name: c.team.String(), Headers: defaultHeaders})
	}

	return c.runTransaction(ctx, tx)
}

type Factory interface {
	New(ctx context.Context, teamName string) (*Client, error)
}

type ClientFactory struct{}

func (cf *ClientFactory) New(ctx context.Context, name string) (*Client, error) {
	client := New(TeamName(name), WithAppEngine(), WithSlackFormatter())

	return client, client.ensureTeamExists(ctx)
}
