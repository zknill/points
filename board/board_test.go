package board

import (
	"context"
	"fmt"
	"os"
	"testing"

	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

var aectx context.Context

func TestMain(m *testing.M) {
	ctx, done, err := aetest.NewContext()
	checkError(err)
	defer done()
	aectx = ctx

	resetTestData()

	os.Exit(m.Run())
}

func resetTestData() {
	storeEntry("alice", 3)
	storeEntry("bob", 2)
	storeEntry("jane", 1)
}

func storeEntry(name string, points int) {
	_, err := datastore.Put(aectx, entryKey(name), &storedEntry{name, points})
	checkError(err)
}

func entryKey(name string) *datastore.Key {
	return datastore.NewKey(aectx, kindEntry, name, 0, nil)
}

func checkError(err error) {
	if err != nil {
		panic(fmt.Sprintf("failed setting up app engine, %+v", err))
	}
}
