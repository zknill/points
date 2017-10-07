package points

import (
	"io"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/zknill/points/board"
	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
)

func TestMsg_String(t *testing.T) {
	tests := []struct {
		name string
		msg  Message
		want string
	}{
		{"msg", &msg{"borked!"}, "borked!"},
		{"add", new(addCmd), addIncorrectTokensText},
		{"list", new(listCmd), "list uses the format: `/points list`"},
		{"reset", new(resetCmd), "reset uses the format: `/points reset`"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.String(); got != tt.want {
				t.Errorf("msg.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddCmd_Execute(t *testing.T) {
	tests := []struct {
		name      string
		cmd       *addCmd
		want      string
		wantErr   bool
		addFunc   func(context.Context, int) error
		addCalled bool
	}{
		{
			name: "success",
			cmd:  &addCmd{tokens: []string{"add", "slackbot"}},
			addFunc: func(_ context.Context, num int) error {
				return nil
			},
			want: "added a point to slackbot",
		},
		{
			name: "bad tokens",
			cmd:  &addCmd{tokens: []string{"add"}},
			addFunc: func(_ context.Context, num int) error {
				return nil
			},
			wantErr:   true,
			addCalled: true,
		},
		{
			name: "bad add func",
			cmd:  &addCmd{tokens: []string{"add", "slackbot"}},
			addFunc: func(_ context.Context, num int) error {
				return errors.New("add func returned error")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			standings := &mockStandings{store: map[string]board.Entry{"slackbot": &mockEntry{
				add: func(ctx context.Context, num int) error {
					tt.addCalled = !tt.addCalled
					return tt.addFunc(ctx, num)
				},
			}}}
			got, err := tt.cmd.Execute(aectx, standings)
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr = %v, err = %v", tt.wantErr, err)
			}
			if got.String() != tt.want {
				t.Errorf("want = %s, got = %s", tt.want, got)
			}
			if !tt.addCalled {
				t.Errorf("incorrect tt.addCalled state, %v", tt.addCalled)
			}
		})
	}
}

func TestListExecute(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() (board.Leaderboard, error)
		want    string
		wantErr bool
	}{
		{
			name: "empty",
			setup: func() (board.Leaderboard, error) {
				return board.Load(aectx, board.NewTeam("empty"))
			},
			want: "+------+--------+\n| NAME | POINTS |\n+------+--------+\n+------+--------+\n",
		},
		{
			name: "one entry",
			setup: func() (board.Leaderboard, error) {
				lb, err := board.Load(aectx, board.NewTeam("one-entry"))
				if err != nil {
					return nil, err
				}
				err = lb.Entry(aectx, "slackbot").Add(aectx, 1)
				return lb, err
			},
			want: "+----------+--------+\n|   NAME   | POINTS |\n+----------+--------+\n| slackbot |      1 |\n+----------+--------+\n",
		},
		{
			name: "no leaderboard key",
			setup: func() (board.Leaderboard, error) {
				return nil, nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		lb, err := tt.setup()
		if err != nil {
			t.Errorf("%q. failed in setup: %v", tt.name, err)
		}
		got, err := new(listCmd).Execute(aectx, lb)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. wantErr = %t, got = %v", tt.name, tt.wantErr, err)
		}
		if got.String() != tt.want {
			t.Errorf("%q. want = %s, got = %s", tt.name, tt.want, got)
		}
	}
}

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

type mockStandings struct {
	store map[string]board.Entry
}

func (s *mockStandings) Entry(_ context.Context, name string) board.Entry {
	return s.store[name]
}

func (s *mockStandings) Out(_ context.Context, w io.Writer) error {
	panic("implement me")
}

func (s *mockStandings) Reset(_ context.Context) error {
	panic("implement me")
}

var _ board.Transacter = new(mockTransacter)

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
