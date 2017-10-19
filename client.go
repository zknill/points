package points

import (
	"context"
	"io"

	"github.com/zknill/points/internal/log"
	"github.com/zknill/points/internal/store"
	"google.golang.org/appengine/datastore"
)

var defaultHeaders = []string{"name", "points"}

// EntryName is the name of an entry
type EntryName string

// String returns teh string representation of an EntryName
func (n EntryName) String() string {
	return string(n)
}

// Entry represents the a single entry
type Entry struct {
	Name  string
	Score int
}

// TeamName is the name of a team
type TeamName string

// String returns the string representation of a TeamName
func (t TeamName) String() string {
	return string(t)
}

// Team represents a team. With the headers for that teams leaderboard.
type Team struct {
	Name    string
	Headers []string
}

type entryStore interface {
	Entry(ctx context.Context, team, name string) (*store.Entry, error)
	Entries(ctx context.Context, team string) ([]*store.Entry, error)
	PutEntry(ctx context.Context, team string, entry *store.Entry) error

	Team(ctx context.Context, team string) (*store.Team, error)
	PutTeam(ctx context.Context, team *store.Team) error
}

type transacter interface {
	Transact(ctx context.Context, tx store.Transaction) error
}

var _ entryStore = (*store.AppEngine)(nil)
var _ transacter = (*store.AppEngine)(nil)

// ClientOption is a func that adapts a client returning
// a new client. It can be used to create defaults that
// set up a *TeamClient with a set of fields.
type ClientOption func(client *TeamClient) *TeamClient

// TeamClient is a object that controls points for a single team.
// It has a number of methods for getting and modifying the
// state of entries for that team. All exported methods of
// a team client should run in transactions.
type TeamClient struct {
	team      TeamName
	store     entryStore
	formatter Formatter
	log       log.Logger
}

// New is a constructor for a *TeamClient. It applies all the
// ClientOptions to the client in the construction.
// ClientOptions should be used to modify the fields of a client.
// ClientOptions take precedence in reverse order,
// the last option to be passed will overwrite any others it
// conflicts with.
func New(team TeamName, options ...ClientOption) *TeamClient {
	client := &TeamClient{
		team: team,
		log:  log.DefaultLogger(),
	}

	for _, opt := range options {
		client = opt(client)
	}

	return client
}

// WithAppEngine returns a ClientOption that sets up a
// TeamClient to work well with google app engine.
func WithAppEngine() ClientOption {
	return func(client *TeamClient) *TeamClient {
		client.log = &log.AppEngine{}
		client.store = &store.AppEngine{}

		return client
	}
}

// WithSlackFormatter returns a ClientOption that
// modifies a client to format it's output in a slack
// friendly way.
func WithSlackFormatter() ClientOption {
	return func(client *TeamClient) *TeamClient {
		client.formatter = formatFunc(slackFormatter)
		return client
	}
}

// Entry the entry with the given entry name. Entry uses
// the team name from the client to get the entry. It returns a
// nil error and an *Entry if the entry was found. A non-nil
// error otherwise.
func (c *TeamClient) Entry(ctx context.Context, name EntryName) (*Entry, error) {
	var entry = new(Entry)
	var tx store.Transaction = func(tc context.Context) error {
		var err error
		entry, err = c.getEntry(tc, name)
		return err
	}

	return entry, c.runTransaction(ctx, tx)
}

// Scores writes to the writer w the current state of the teams entries.
// It uses the team name and formatter from the *TeamClient to get
// the entries and format the output. A non-nil error will be returned
// if Scores was not successful. Scores is transactional.
func (c *TeamClient) Scores(ctx context.Context, w io.Writer) error {
	var err error
	var team *store.Team
	var entries []*Entry

	var tx store.Transaction = func(tc context.Context) error {
		team, err = c.store.Team(ctx, c.team.String())
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

// Add increases the score of a single entry by one, given the entryName
// A non-nil error indicates Add failed. Add is a transactional operation.
func (c *TeamClient) Add(ctx context.Context, name EntryName) error {

	var tx store.Transaction = func(tc context.Context) error {

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

// Reset updates all the points for entries of a team to zero.
// It is a transactional operation, a non-nil error indicates failure.
func (c *TeamClient) Reset(ctx context.Context) error {

	var tx store.Transaction = func(tc context.Context) error {

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

func (c *TeamClient) getEntry(ctx context.Context, name EntryName) (*Entry, error) {
	storeEntry, err := c.store.Entry(ctx, c.team.String(), name.String())
	if err != nil {
		return nil, err
	}
	e := Entry(*storeEntry)
	return &e, nil
}

func (c *TeamClient) putEntry(ctx context.Context, entry *Entry) error {
	s := store.Entry(*entry)
	return c.store.PutEntry(ctx, c.team.String(), &s)
}

func (c *TeamClient) getAll(ctx context.Context) ([]*Entry, error) {
	sEntries, err := c.store.Entries(ctx, c.team.String())

	entries := make([]*Entry, len(sEntries))
	for i := range sEntries {
		e := Entry(*sEntries[i])
		entries[i] = &e
	}

	return entries, err
}

func (c *TeamClient) runTransaction(ctx context.Context, tx store.Transaction) error {
	if d, ok := c.store.(transacter); ok {
		return d.Transact(ctx, tx)
	}
	return tx(ctx)
}

func (c *TeamClient) ensureTeamExists(ctx context.Context) error {
	var tx store.Transaction = func(tc context.Context) error {
		_, err := c.store.Team(ctx, c.team.String())
		if err == nil {
			return nil
		}

		if err != datastore.ErrNoSuchEntity {
			return err
		}

		return c.store.PutTeam(ctx, &store.Team{Name: c.team.String(), Headers: defaultHeaders})
	}

	return c.runTransaction(ctx, tx)
}

// Factory defines methods for creating a *TeamClient.
// Factories should set up *TeamClient for operating
// with the correct loggers, formatters, datastores etc.
type Factory interface {
	New(ctx context.Context, teamName string) (*TeamClient, error)
}

// AppEngineFactory is a Factory for creating TeamClients
// that work well with google app engine.
type AppEngineFactory struct{}

// New implements the Factory interface for AppEngineFactory.
// It creates a new *TeamClient, ensures that the team exists
// and sets up the datastore, logger and formatter for appengine.
func (cf *AppEngineFactory) New(ctx context.Context, name string) (*TeamClient, error) {
	client := New(TeamName(name), WithAppEngine(), WithSlackFormatter())

	return client, client.ensureTeamExists(ctx)
}
