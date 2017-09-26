package board

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const kindEntry string = "Entries"

// ErrEntryNotFound represents and entry that could not be found
type ErrEntryNotFound struct {
	// Name is the name of the entry that was not found
	Name string
}

// Error implements the Error interface for ErrEntryNotFound
func (e *ErrEntryNotFound) Error() string {
	return fmt.Sprintf("not found: %s", e.Name)
}

// ErrBadKey represents an error that's returned because the key used
// for the lookup or put was incorrect because entity was not found.
type ErrBadKey string

// Error implements the Error interface for ErrBadKey
func (e ErrBadKey) Error() string {
	return fmt.Sprintf("bad key: %s", string(e))
}

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
	// Delete removes the current entry from the datastore.
	Delete(ctx context.Context) error
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
		if _, ok := errors.Cause(err).(*ErrEntryNotFound); ok {
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

// Delete removes the current entry from the datastore
func (aee *aeEntry) Delete(ctx context.Context) error {
	return aee.del(ctx)
}

func (aee *aeEntry) get(ctx context.Context) (*storedEntry, error) {
	log.Debugf(ctx, "attempt to get entry for key: %s", aee.entryKey)
	var stored = new(storedEntry)
	if err := datastore.Get(ctx, aee.entryKey, stored); err != nil {
		if err == datastore.ErrNoSuchEntity {
			log.Debugf(ctx, "failed to get key %s", aee.entryKey)
			return nil, errors.Wrap(&ErrEntryNotFound{Name: aee.entryKey.StringID()}, "failed to get")
		} else if err == datastore.ErrNoSuchEntity || err == datastore.ErrInvalidEntityType {
			return nil, errors.Wrapf(ErrBadKey(err.Error()), "failed to get key: %s", aee.entryKey)
		}
	}
	log.Infof(ctx, "got entry for key: %s", aee.entryKey)
	return stored, nil
}

func (aee *aeEntry) put(ctx context.Context, entry *storedEntry) error {
	log.Debugf(ctx, "attempt to store entry for key: %s, stored entry: %+v", aee.entryKey, entry)
	if _, err := datastore.Put(ctx, aee.entryKey, entry); err != nil {
		log.Infof(ctx, "failed to store entry for key: %s", aee.entryKey)
		if err == datastore.ErrNoSuchEntity || err == datastore.ErrInvalidEntityType {
			return errors.Wrapf(ErrBadKey(err.Error()), "failed to put key: %s", aee.entryKey)
		}
		return errors.Wrapf(err, "failed to stored entry for key: %s", aee.entryKey)
	}
	log.Infof(ctx, "stored entry for key: %s, stored entry: %+v", aee.entryKey, entry)
	return nil
}

func (aee *aeEntry) del(ctx context.Context) error {
	log.Debugf(ctx, "attempt to delete entry for key: %s", aee.entryKey)
	if err := datastore.Delete(ctx, aee.entryKey); err != nil {
		log.Infof(ctx, "failed to delete entry for key: %s", aee.entryKey)
		if err == datastore.ErrNoSuchEntity {
			return errors.Wrapf(ErrBadKey(err.Error()), "failed to delete key: %s", aee.entryKey)
		}
		return errors.Wrapf(err, "failed to delete key: %s", aee.entryKey)
	}
	return nil
}
