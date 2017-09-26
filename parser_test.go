package points

import (
	"io"
	"reflect"
	"testing"

	"errors"

	"os"

	"github.com/zknill/points/board"
	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
)

var aectx context.Context

func TestMain(m *testing.M) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		panic(err)
	}
	defer done()
	aectx = ctx

	os.Exit(m.Run())
}

func Test_add_success(t *testing.T) {
	tokens := []string{"add", "slackbot"}

	ctx := context.Background()
	points := 1

	var addCalled bool
	addFunc := func(_ context.Context, num int) error {
		if num != points {
			t.Errorf("incorrect points, wanted %v got %v", points, num)
		}
		addCalled = true
		return nil
	}

	entry := &mockEntry{add: addFunc}
	standings := &mockStandings{store: map[string]board.Entry{"slackbot": entry}}

	msg := add(tokens)(ctx, standings)
	t.Log(msg)

	if !addCalled {
		t.Error("expected add func to be called")
	}
}

func Test_add_success_transaction(t *testing.T) {
	tokens := []string{"add", "slackbot"}

	ctx := context.Background()
	points := 1

	var addCalled bool
	addFunc := func(_ context.Context, num int) error {
		if num != points {
			t.Errorf("incorrect points, wanted %v got %v", points, num)
		}
		addCalled = true
		return nil
	}

	entry := &mockEntry{add: addFunc}
	standings := &mockTransacter{mockStandings: mockStandings{store: map[string]board.Entry{"slackbot": entry}}}

	msg := add(tokens)(ctx, standings)
	t.Log(msg)

	if !addCalled {
		t.Error("expected add func to be called")
	}

	if !standings.called {
		t.Error("expected add to be run in transaction")
	}
}

func Test_add_error(t *testing.T) {
	tokens := []string{"add", "slackbot"}

	points := 1

	var addCalled bool
	addFunc := func(_ context.Context, num int) error {
		if num != points {
			t.Errorf("incorrect points, wanted %v got %v", points, num)
		}
		addCalled = true
		return errors.New("something's broken")
	}

	entry := &mockEntry{add: addFunc}
	standings := &mockStandings{store: map[string]board.Entry{"slackbot": entry}}

	msg := add(tokens)(aectx, standings)
	t.Log(msg)

	if !addCalled {
		t.Error("expected add func to be called")
	}
}

func Test_errorMessage(t *testing.T) {
	tests := []struct {
		name string
		err  string
	}{
		{"simple", "error"},
		{"spaces", "error message"},
		{"new line", "error\nmessage"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errorMessage(tt.err)(nil, nil); !reflect.DeepEqual(got, tt.err) {
				t.Errorf("errorMessage() = %v, want %v", got, tt.err)
			}
		})
	}
}

func Test_helpMessage(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"simple", helpText},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helpMessage()(nil, nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("helpMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockStandings struct {
	store map[string]board.Entry
}

func (s *mockStandings) Entry(_ context.Context, name string) board.Entry {
	return s.store[name]
}

func (s *mockStandings) Out(_ context.Context, w io.Writer) {
	panic("implement me")
}

func (s *mockStandings) Reset(_ context.Context) error {
	panic("implement me")
}

type mockTransacter struct {
	mockStandings
	called bool
}

func (t *mockTransacter) Transact(ctx context.Context, f func(tc context.Context) error) error {
	t.called = true
	return f(ctx)
}

type mockEntry struct {
	add   func(ctx context.Context, num int) error
	score func(ctx context.Context) (int, error)
}

func (e *mockEntry) Add(ctx context.Context, num int) error {
	return e.add(ctx, num)
}

func (e *mockEntry) Score(ctx context.Context) (int, error) {
	return e.score(ctx)
}

func (e *mockEntry) Delete(ctx context.Context) error {
	panic("implement me!")
}
