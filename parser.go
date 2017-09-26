package points

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/zknill/points/board"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

const (
	helpText                  string = "try one of: help, list, add, reset"
	addIncorrectTokensText    string = "add needs the format `/points add slackbot`"
	renameIncorrectTokensText string = "rename needs the format `/points rename slackbot pointsbot`"
	createErrorText           string = "create a team using `/points init`"
)

// Command is the function to be called to handle a users request.
// It returns a string that can be sent back to the user.
type Command func(ctx context.Context, standings board.Leaderboard) string

// Parser takes the raw string text for the message written by the user and returns a command.
// Some returned commands are wrapped in validation middleware to ensure they are passed legitimate arguments.
func Parser(ctx context.Context, team, request string) Command {
	tokens := strings.Fields(strings.TrimSpace(request))

	log.Debugf(ctx, "parsing tokens: %v, raw: %v, team: %q", tokens, request, team)

	if len(tokens) < 1 {
		return helpMessage()
	}

	switch tokens[0] {
	case "init":
		return initTeam(team)
	case "add":
		return validateStandings(add(tokens))
	case "rename":
		return validateStandings(rename(tokens))
	case "reset":
		return validateStandings(reset())
	case "list":
		return validateStandings(list())
	case "help":
		return helpMessage()
	default:
		return errorMessage(fmt.Sprintf("unrecognised command %q", tokens[0]))
	}
}

func add(tokens []string) Command {
	if len(tokens) != 2 {
		return errorMessage(addIncorrectTokensText)
	}
	name := strings.ToLower(tokens[1])
	return func(ctx context.Context, standings board.Leaderboard) string {

		addFunc := func(ctx context.Context) error {
			return standings.Entry(ctx, name).Add(ctx, 1)
		}

		if ts, ok := standings.(board.Transacter); ok {
			if err := ts.Transact(ctx, addFunc); err != nil {
				log.Errorf(ctx, "failed adding points to %q: %+s", name, err)
				return response(err)
			}
		} else {
			if err := addFunc(ctx); err != nil {
				log.Errorf(ctx, "failed adding points to %q: %+s", name, err)
				return response(err)
			}
		}
		return fmt.Sprintf("added a point to %q", name)
	}
}

func initTeam(name string) Command {
	return func(ctx context.Context, standings board.Leaderboard) string {

		createFunc := func(c context.Context) error {
			return board.Create(c, board.NewTeam(name))
		}

		if ts, ok := standings.(board.Transacter); ok {
			if err := ts.Transact(ctx, createFunc); err != nil {
				log.Errorf(ctx, "failed creating standings for team %q", name)
				return response(err)
			}
		} else {
			if err := createFunc(ctx); err != nil {
				log.Errorf(ctx, "failed creating standings for team %q", name)
				return response(err)
			}
		}

		return "init success!"
	}
}

func list() Command {
	return func(ctx context.Context, standings board.Leaderboard) string {
		var buf bytes.Buffer
		standings.Out(ctx, &buf)
		return buf.String()
	}
}

func reset() Command {
	return func(ctx context.Context, standings board.Leaderboard) string {
		resetFunc := func(ctx context.Context) error {
			return standings.Reset(ctx)
		}

		if ts, ok := standings.(board.Transacter); ok {
			if err := ts.Transact(ctx, resetFunc); err != nil {
				return response(err)
			}
		} else if err := resetFunc(ctx); err != nil {
			return response(err)
		}

		return "reset all points!"
	}
}

func rename(tokens []string) Command {
	return func(ctx context.Context, standings board.Leaderboard) string {
		if len(tokens) != 3 { // []string{rename, old_name, new_name}
			return renameIncorrectTokensText
		}
		oldName, newName := tokens[1], tokens[2]

		renameFunc := func(c context.Context) error {
			old := standings.Entry(ctx, oldName)

			if _, err := standings.Entry(ctx, newName).Score(ctx); err == nil {
				return errors.New("%q does not exist")
			} else if _, ok := errors.Cause(err).(*board.ErrEntryNotFound); !ok {
				return errors.Wrap(err, "unexpected failure")
			}

			score, err := old.Score(ctx)
			if err != nil {
				return errors.Wrap(err, "failed getting original score")
			}

			if delErr := old.Delete(ctx); delErr != nil {
				return errors.Wrap(delErr, "failed deleting existing entry")
			}

			if err := standings.Entry(ctx, newName).Add(ctx, score); err != nil {
				return errors.Wrap(err, "failed creating new entry")
			}

			return nil
		}

		if ts, ok := standings.(board.Transacter); ok {
			if err := ts.Transact(ctx, renameFunc); err != nil {
				log.Errorf(ctx, "failed renaming %q to %q: %+s", oldName, newName, err)
				return response(err)
			}
		}
		if err := renameFunc(ctx); err != nil {
			log.Errorf(ctx, "failed renaming %q to %q: %+s", oldName, newName, err)
			return response(err)
		}
		return fmt.Sprintf("renamed %q to %q", oldName, newName)
	}
}

func errorMessage(err string) Command {
	return func(_ context.Context, _ board.Leaderboard) string {
		return err
	}
}

func helpMessage() Command {
	return func(_ context.Context, _ board.Leaderboard) string {
		return helpText
	}
}

func validateStandings(command Command) Command {
	return func(ctx context.Context, standings board.Leaderboard) string {
		if standings == nil {
			return createErrorText
		}
		return command(ctx, standings)
	}
}

func response(err error) string {
	if err == nil {
		return "whoops! something failed"
	}
	switch e := errors.Cause(err).(type) {
	case *board.ErrEntryNotFound:
		return fmt.Sprintf("could not find entry %q", e.Name)
	default:
		return "ack! something went wrong..."
	}
}
