package points

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"google.golang.org/appengine/log"
)

const errorMessage = "ack! something went wrong"

var Controller = &SlackController{
	factory: ClientFactory{},
}

type SlackController struct {
	factory ClientFactory
}

func (c *SlackController) Parse(ctx context.Context, team string, request string) string {

	tokens := strings.Fields(request)

	client, err := c.factory.New(ctx, team)
	if err != nil {
		log.Errorf(ctx, "factory method failed creating new client, %+v", err)
		return errorMessage
	}

	switch tokens[0] {
	case "add":
		return add(ctx, client, tokens)
	case "list":
		return list(ctx, client)
	case "reset":
		return reset(ctx, client)
	}

	return fmt.Sprintf("unknown command %q", tokens[0])
}

func add(ctx context.Context, client *Client, tokens []string) string {

	// check tokens length

	name := EntryName(tokens[1])

	err := client.Add(ctx, name)
	if err != nil {
		return errorMessage
	}

	return fmt.Sprintf("added a point to %q", name.String())
}

func list(ctx context.Context, client *Client) string {
	var buf bytes.Buffer
	if err := client.Scores(ctx, &buf); err != nil {
		return errorMessage
	}

	return buf.String()
}

func reset(ctx context.Context, client *Client) string {
	if err := client.Reset(ctx); err != nil {
		return errorMessage
	}
	return "reset all points"
}
