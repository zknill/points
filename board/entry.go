package board

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const kindEntry string = "Entries"

// ErrEntryNotFound represents and entry that could not be found.
var ErrEntryNotFound error = errors.New("entry: not found")

// Entry represents a single row in a Leaderboard
type Entry interface {
	// Add increases the score that an entry has by num provided.
	// If Entry is backed by a datastore, you can expect that after Add returns, the store will be updated.
	// Add should be wrapped in a transaction where applicable.
	Add(ctx context.Context, num int) error
	// Score returns the number of points for an entry. There's no guarantee when the score is retrieved.
	// If the entry is stored in a datastore then,the number of points could be retrieved at the point when the Entry
	// is retrieved from the datastore, or when Score is called.
	// This is less of a problem if Score is called in a transaction.
	Score(ctx context.Context) (int, error)
}

var _ Entry = (*aeEntry)(nil)

// aeEntry is an entry that stores information for how to interact with the entry in the datastore.
// It does not contain the information from the entry in the datastore.
// That information is accessed and modified in the methods that have an aeEntry receiver.
type aeEntry struct {
	entryKey *datastore.Key
}

// storedEntry represents the information that's actually put into the datastore.
// It's used by aeEntry's get and put methods.
type storedEntry struct {
	Name   string
	Points int
}

// Add a point to an entry stored in app engine datastore.
// If the entry is not found, Add creates it with the number of points given.
func (aee *aeEntry) Add(ctx context.Context, num int) error {
	stored, err := aee.get(ctx)
	if err != nil {
		if errors.Cause(err) == ErrEntryNotFound {
			stored = &storedEntry{Name: aee.entryKey.StringID()}
		} else {
			return errors.Wrapf(err, "failed to add %d points", num)
		}
	}
	stored.Points += num
	return aee.put(ctx, stored)
}

// Score gets the score of an entry. The points are retrieved when score is called.
func (aee *aeEntry) Score(ctx context.Context) (int, error) {
	entry, err := aee.get(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed getting score")
	}
	return entry.Points, nil
}

func (aee *aeEntry) get(ctx context.Context) (*storedEntry, error) {
	log.Debugf(ctx, "attempt to get entry for key: %s", aee.entryKey)
	var stored = new(storedEntry)
	if err := datastore.Get(ctx, aee.entryKey, stored); err != nil {
		if err == datastore.ErrNoSuchEntity {
			log.Debugf(ctx, "failed to get key %s", aee.entryKey)
			return nil, errors.Wrapf(ErrEntryNotFound, "failed to get key: %s", aee.entryKey)
		}
	}
	log.Infof(ctx, "got entry for key: %s", aee.entryKey)
	return stored, nil
}

func (aee *aeEntry) put(ctx context.Context, entry *storedEntry) error {
	log.Debugf(ctx, "attempt to store entry for key: %s, stored entry: %+v", aee.entryKey, entry)
	if _, err := datastore.Put(ctx, aee.entryKey, entry); err != nil {
		log.Infof(ctx, "failed to store entry for key: %s", aee.entryKey)
		return errors.Wrapf(err, "failed to stored entry for key: %s", aee.entryKey)
	}
	log.Infof(ctx, "stored entry for key: %s, stored entry: %+v", aee.entryKey, entry)
	return nil
}
