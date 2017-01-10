package points

import (
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/zknill/points/commands"
	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

func Test_storeEntry_getEntry(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer done()

	type args struct {
		ctx   context.Context
		entry *points.Entry
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "simple entry",
			args: args{
				ctx: ctx,
				entry: &points.Entry{
					Name:   "Gordon",
					Points: 2,
				},
			},
			wantErr: false,
		},
		{
			name: "another simple entry",
			args: args{
				ctx: ctx,
				entry: &points.Entry{
					Name:   "Jane",
					Points: 7,
				},
			},
			wantErr: false,
		},
		{
			name: "store nil entry",
			args: args{
				ctx:   ctx,
				entry: nil,
			},
			wantErr: true,
		},
		{
			name: "store entry with no name",
			args: args{
				ctx: ctx,
				entry: &points.Entry{
					Name:   "",
					Points: 9,
				},
			},
			wantErr: true,
		},
		{
			name: "store entry with no points",
			args: args{
				ctx: ctx,
				entry: &points.Entry{
					Name: "Neil",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if err := storeEntry(tt.args.ctx, tt.args.entry); (err != nil) != tt.wantErr {
			t.Errorf("%q. storeEntry() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		if tt.args.entry != nil {
			entry, err := getEntry(ctx, strings.ToLower(tt.args.entry.Name))
			if (err != nil) != tt.wantErr {
				t.Errorf("%q. getEntry() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				continue
			}
			if !reflect.DeepEqual(entry, tt.args.entry) && !tt.wantErr {
				t.Errorf("%q. getEntry() = %v, want %v", tt.name, entry, tt.args.entry)
			}
		}
	}
}

func Test_getEntries(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer done()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *[]*points.Entry
		wantErr bool
	}{
		{
			name: "simple entries",
			args: args{ctx},
			want: &[]*points.Entry{
				{
					Name:   "Neil",
					Points: 3,
				},
				{
					Name:   "Jane",
					Points: 2,
				},
				{
					Name:   "Gordon",
					Points: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		for _, ent := range *tt.want {
			if err := storeEntry(tt.args.ctx, ent); err != nil {
				log.Fatalf("%q. failed to store entry = %s in setup", tt.name, ent)
			}
			if _, err := getEntry(tt.args.ctx, ent.Name); err != nil {
				log.Fatalf("%q. entry = %s was not in storage", tt.name, ent)
			}
		}
		got, err := getEntries(tt.args.ctx)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. getEntries() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. getEntries() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_getLeaderboard(t *testing.T) {
	_, done, err := aetest.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer done()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *points.StoredLeaderboard
		wantErr bool
	}{
	// TODO: Add test cases.
	//	{
	//		name: "simple get leaderboard",
	//		args: args{ctx},
	//		want: &points.StoredLeaderboard{
	//			Headers: []string{"header1", "header2"},
	//		},
	//		wantErr: false,
	//	},
	}
	for _, tt := range tests {
		initLeaderboard(tt.args.ctx, tt.want.Headers)
		got, err := getLeaderboard(tt.args.ctx)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. getLeaderboard() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. getLeaderboard() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_initLeaderboard(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer done()

	if lb, err := getLeaderboard(ctx); err == nil {
		log.Fatalf("leaderboard should not exist, leaderboard: %s", lb)
	}

	type args struct {
		ctx     context.Context
		headers []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "simple leaderboard init",
			args: args{
				ctx:     ctx,
				headers: []string{"name", "points"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if err := initLeaderboard(tt.args.ctx, tt.args.headers); (err != nil) != tt.wantErr {
			t.Errorf("%q. initLeaderboard() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_entryKey(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer done()

	type args struct {
		ctx   context.Context
		entry string
	}
	tests := []struct {
		name string
		args args
		want *datastore.Key
	}{
		{
			name: "simple key",
			args: args{
				ctx,
				"gordon",
			},
			want: datastore.NewKey(ctx, ENTRY, "Gordon", 0, nil),
		},
	}
	for _, tt := range tests {
		if got := entryKey(tt.args.ctx, tt.args.entry); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. entryKey() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
