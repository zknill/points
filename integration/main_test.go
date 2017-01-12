package points

import (
	"log"
	"os"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
)

var ctx context.Context

func TestMain(m *testing.M) {
	var err error
	var done func()
	ctx, done, err = aetest.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer done()
	code := m.Run()
	ctx.Done()
	os.Exit(code)
}
