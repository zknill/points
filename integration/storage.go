package points

import (
	"errors"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/zknill/points/commands"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func storeEntry(ctx context.Context, entry *points.Entry) error {
	log.Infof(ctx, "Attempt to put entries into storage")
	h := hash(entry.Name)
	log.Infof(ctx, fmt.Sprintf("entry name %s hash %v", entry.Name, h))
	key := datastore.NewKey(ctx, ENTRY, entry.Name, 0, nil)

	if _, err := datastore.Put(ctx, key, entry); err != nil {
		message := fmt.Sprintf("put entry in storage failed, error %s", err.Error())
		log.Errorf(ctx, message)
		return errors.New(message)
	}
	log.Infof(ctx, "Successfully put entry in storage")
	return nil
}

func hash(s string) int64 {
	h := fnv.New32a()
	h.Write([]byte(s))
	n := int32(h.Sum32())
	if n < 0 {
		return int64(^uint32(n - 1))
	}
	return int64(n)
}

func getEntries(ctx context.Context) (*[]*points.Entry, error) {
	log.Infof(ctx, "attempt to get entry from storage")
	entries := new([]*points.Entry)

	q := datastore.NewQuery(ENTRY).Order("-Points")

	if _, err := q.GetAll(ctx, entries); err != nil {
		message := "get all entries failed. Error: " + err.Error()
		log.Errorf(ctx, message)
		return nil, errors.New(message)
	}
	log.Infof(ctx, fmt.Sprintf("successfully retrieved %d entries from storage", len(*entries)))
	return entries, nil
}

func getLeaderboard(ctx context.Context) (*points.StoredLeaderboard, error) {
	log.Infof(ctx, "attempt to get stored leaderboard from storage")
	lbSlice := new([]*points.StoredLeaderboard)

	q := datastore.NewQuery(LEADERBOARD).Limit(1)

	if _, err := q.GetAll(ctx, lbSlice); err != nil {
		message := "failed to get stored leaderboard from storage. Error: " + err.Error()
		log.Errorf(ctx, message)
		return nil, errors.New(message)
	}

	if len(*lbSlice) != 0 {
		return (*lbSlice)[0], nil
	}
	return nil, errors.New("Leaderboard does not exist")
}

func initLeaderboard(ctx context.Context, headers []string) error {
	log.Infof(ctx, "attempt to init new leaderboard with headers: ", strings.Join(headers, ", "))

	storeLb := &points.StoredLeaderboard{
		Headers: headers,
	}

	key := datastore.NewIncompleteKey(ctx, LEADERBOARD, nil)

	if _, err := datastore.Put(ctx, key, storeLb); err != nil {
		message := "failed to init new leaderboard with headers: " + strings.Join(headers, ", ") + "Error: " + err.Error()
		log.Errorf(ctx, message)
		return errors.New(message)
	}
	return nil
}
