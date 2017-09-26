package board

import (
	"io"
	"sort"
	"strconv"
	"strings"

	"fmt"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const kindLeaderboard string = "Leaderboard"

var defaultHeaders = []string{"name", "points"}

// Leaderboard holds all the points information like a conventional leaderboard.
type Leaderboard interface {
	// Entry gets an entry for the given name
	Entry(ctx context.Context, name string) Entry
	// Out gets all the Entrys in the Leaderboard
	// and writes their names and points to w as an ascii table
	Out(ctx context.Context, w io.Writer)
	// Reset updates all the Entrys in the Leaderboard to have a zero score
	// Entry names are maintained
	Reset(ctx context.Context) error
}

// Transacter defines a transaction that allows the given func to be run in an atomic transaction.
type Transacter interface {
	// Transact takes a function to be called transactionally.
	// The transaction is rolled back if the err returned from the t func is non-nil.
	Transact(ctx context.Context, t func(tc context.Context) error) error
}

// ErrBoardNotFound represents a not found error for leaderboards
type ErrBoardNotFound string

// Error implements the Error interface for ErrBoardNotFound
func (e ErrBoardNotFound) Error() string {
	return fmt.Sprintf("board not found: %s", string(e))
}

// Load returns a Leaderboard for a team
// It uses the appengine datastore as a storage layer.
func Load(ctx context.Context, team Team) (Leaderboard, error) {
	k := team.key(ctx)
	s := new(aeLeaderBoard)
	log.Infof(ctx, "attempt to get standings for key: %s", k)
	if err := datastore.Get(ctx, k, s); err != nil {
		if err == datastore.ErrNoSuchEntity {
			e := ErrBoardNotFound(fmt.Sprintf("team: %s", team.String()))
			return nil, errors.Wrap(e, "failed getting leaderboard")
		}
		return nil, errors.Wrap(err, "failed getting leaderboard")
	}
	s.StandingsKey = k
	log.Infof(ctx, "successfully got standings for key: %s", k)
	return s, nil
}

// Create will stored a new leaderboard for a team with default headers.
// If the team already has a leaderboard, Create will error.
// It uses the appengine datastore as a storage layer.
func Create(ctx context.Context, team Team) error {
	k := team.key(ctx)
	_, err := Load(ctx, team)
	if err != nil && errors.Cause(err) != datastore.ErrNoSuchEntity {
		return errors.Wrapf(err, "failed checking if standings already exist for key: %s", k)
	}
	s := &aeLeaderBoard{
		Team:    team.String(),
		Headers: defaultHeaders,
	}
	log.Infof(ctx, "creating standings for key: %s", k)
	if _, err := datastore.Put(ctx, k, s); err != nil {
		return errors.Wrap(err, "failed storing new standings")
	}
	log.Infof(ctx, "successfully created new standings for key: %s", k)
	return nil
}

// Team defines a slack team
type Team string

// NewTeam is a helper method to create a new Team type
func NewTeam(name string) Team {
	return Team(name)
}

// String implements the Stringer interface and is a helper to get the string out of a team.
func (t Team) String() string {
	return string(t)
}

func (t Team) key(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, kindLeaderboard, t.String(), 0, nil)
}

var _ Leaderboard = (*aeLeaderBoard)(nil)
var _ Transacter = (*aeLeaderBoard)(nil)

type aeLeaderBoard struct {
	StandingsKey *datastore.Key `datastore:"-"`
	Team         string
	Headers      []string
}

// Transact implements the Transactor interface for *aeLeaderboard
func (i *aeLeaderBoard) Transact(ctx context.Context, t func(tc context.Context) error) error {
	return datastore.RunInTransaction(ctx, t, nil)
}

// Entry implements the Leaderboard interface and sets up the information to retrieve an entry.
// Note that the entries information is not accessed until a method is called on that entry.
func (i *aeLeaderBoard) Entry(ctx context.Context, name string) Entry {
	key := datastore.NewKey(ctx, kindEntry, strings.ToLower(name), 0, i.StandingsKey)
	return &aeEntry{entryKey: key}
}

// Out implements Leaderboard and writes all the entries for this leaderboard in a slack friendly way.
func (i *aeLeaderBoard) Out(ctx context.Context, w io.Writer) {
	query := datastore.NewQuery(kindEntry).Ancestor(i.StandingsKey)
	results := query.Run(ctx)
	var entries []*storedEntry
	for {
		entry := new(storedEntry)
		if _, err := results.Next(entry); err != nil {
			if err == datastore.Done {
				break
			}
			log.Errorf(ctx, err.Error())
		}
		entries = append(entries, entry)
	}
	writeTable(w, i.Headers, entries)
}

// Reset implements Leaderboard, it allows all the associated entries to have their points set to zero.
func (i *aeLeaderBoard) Reset(ctx context.Context) error {
	query := datastore.NewQuery(kindEntry).Ancestor(i.StandingsKey)
	results := query.Run(ctx)
	stored := new(storedEntry)
	for {
		k, err := results.Next(stored)
		if err != nil {
			if err == datastore.Done {
				break
			}
			return err
		}
		stored.Points = 0
		entry := &aeEntry{entryKey: k}
		if err := entry.put(ctx, stored); err != nil {
			log.Errorf(ctx, "failed putting entry on reset: %+s", err)
		}
	}
	return nil
}

func writeTable(w io.Writer, headers []string, entries []*storedEntry) {
	sort.Slice(entries, func(i, j int) bool {
		iEntry, jEntry := entries[i], entries[j]
		if iEntry.Points == jEntry.Points {
			return iEntry.Name < jEntry.Name
		}
		return iEntry.Points > jEntry.Points
	})
	tw := tablewriter.NewWriter(w)
	tw.SetHeader(headers)
	for _, entry := range entries {
		tw.Append([]string{entry.Name, strconv.Itoa(entry.Points)})
	}
	tw.Render()
}
