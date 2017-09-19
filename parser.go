package points

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/zknill/points/board"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

const (
	helpText               string = "try one of: help, list, add, reset"
	errorText              string = "ack! something failed."
	addIncorrectTokensText string = "add needs the format `/points add slackbot`"
	addSuccessText         string = "added a point to %s"
	resetSuccessText       string = "reset all points!"
	createErrorText        string = "create a team using `/points init`"
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
		if ts, ok := standings.(board.Transacter); ok {
			transaction := func(tc context.Context) error {
				return standings.Entry(ctx, name).Add(ctx, 1)
			}
			if err := ts.Transact(ctx, transaction); err != nil {
				log.Errorf(ctx, "failed adding: %+s", err)
				return errorText
			}
		} else {
			if err := standings.Entry(ctx, name).Add(ctx, 1); err != nil {
				log.Errorf(ctx, "failed add: %+s", err)
				return errorText
			}
		}
		return fmt.Sprintf(addSuccessText, name)
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
				return errorText
			}
		} else {
			if err := createFunc(ctx); err != nil {
				log.Errorf(ctx, "failed creating standings for team %q", name)
				return errorText
			}
		}

		return fmt.Sprintf("init success!")
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
		standings.Reset(ctx)
		return resetSuccessText
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
