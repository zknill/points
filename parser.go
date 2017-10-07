package points

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

const (
	helpText               string = "try one of: help, list, add, reset"
	errorText              string = "ack! something failed"
	addIncorrectTokensText string = "add needs the format `/points add slackbot`"
	addSuccessText         string = "added a point to %s"
	resetSuccessText       string = "reset all points!"
)

// Parser takes the raw string text for the message written by the user and returns a command.
// Some returned commands are wrapped in validation middleware to ensure they are passed legitimate arguments.
func Parser(ctx context.Context, team, request string) Message {
	tokens := strings.Fields(strings.TrimSpace(request))

	log.Debugf(ctx, "parsing tokens: %v, raw: %v, team: %q", tokens, request, team)

	if len(tokens) < 1 {
		return &msg{helpText}
	}

	switch tokens[0] {
	case "init":
		return &initCmd{teamName: team}
	case "add":
		return &addCmd{tokens: tokens}
	case "reset":
		return new(resetCmd)
	case "list":
		return new(listCmd)
	case "help":
		return &msg{helpText}
	default:
		return &msg{fmt.Sprintf("unrecognised command %q", tokens[0])}
	}
}
