package points

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/zknill/points/board"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

// Message is a basic type that defines a String method
// It is a message to send back to the user. A Message
// is any type that implements fmt.Stringer
type Message interface {
	fmt.Stringer
}

// Command implements the Command pattern and
// defines a single method interface to execute
// the underlying command. Returning the message
// to send back to the user. Commands need not be
// threadsafe. Commands should check if the passed
// board.Leaderboard implements the board.Transacter
// interface and execute their operations inside
// a transaction if the board is a Transacter.
type Command interface {
	Execute(ctx context.Context, standings board.Leaderboard) (response Message, err error)
}

type msg struct {
	msg string
}

// String implements fmt.Stringer for
// msg making it a Message.
func (e *msg) String() string {
	return e.msg
}

func message(text string) Message {
	return &msg{msg: text}
}

type transaction func(tc context.Context) error

type addCmd struct {
	tokens []string
}

// String implements fmt.Stringer for
// addCmd making it a Message.
func (*addCmd) String() string {
	return addIncorrectTokensText
}

// Execute implements the Command interface making
// addCmd a command, it adds a single point to an
// entry based on the tokenised input []string.
func (cmd *addCmd) Execute(ctx context.Context, standings board.Leaderboard) (Message, error) {
	if standings == nil {
		return nil, errors.New("could not find leaderboard")
	}

	if len(cmd.tokens) != 2 {
		return nil, errors.New(addIncorrectTokensText)
	}

	name := strings.ToLower(cmd.tokens[1])

	var addFunc transaction = func(tc context.Context) error {
		return standings.Entry(tc, name).Add(tc, 1)
	}

	addTx := wrapTransacter(standings, addFunc)
	if err := executeTransaction(ctx, addTx); err != nil {
		return nil, err
	}

	return message(fmt.Sprintf(addSuccessText, name)), nil
}

type listCmd struct {}

// String implements fmt.Stringer for
// listCmd making it a Message.
func (*listCmd) String() string {
	return "list uses the format: `/points list`"
}

// Execute implements the Command interface for listCmd
// listCmd's Execute method lists all the entries for a
// board.Leaderboard in a slack friendly way.
func (cmd *listCmd) Execute(ctx context.Context, standings board.Leaderboard) (response Message, err error) {
	if standings == nil {
		return nil, errors.New("could not find leaderboard")
	}

	var buf = new(bytes.Buffer)
	var listFunc transaction = func(tc context.Context) error {
		return standings.Out(tc, buf)
	}

	listTx := wrapTransacter(standings, listFunc)
	if err := executeTransaction(ctx, listTx); err != nil {
		return nil, err
	}

	return buf, nil
}

type initCmd struct {
	teamName string
}

// String implements fmt.Stringer for
// initCmd making it a Message.
func (*initCmd) String() string {
	return "init uses the format: `/points init`"
}

// Execute for initCmd - is the necessary?
func (cmd *initCmd) Execute(ctx context.Context, standings board.Leaderboard) (response Message, err error) {
	if standings == nil {
		return nil, errors.New("could not find leaderboard")
	}
	var initFunc transaction = func(tc context.Context) error {
		return board.Create(tc, board.NewTeam(cmd.teamName))
	}

	initTx := wrapTransacter(standings, initFunc)
	if err := executeTransaction(ctx, initTx); err != nil {
		return nil, err
	}

	return message("init success!"), nil
}

type resetCmd struct {}

// String implements fmt.Stringer for
// resetCmd making it a Message.
func (*resetCmd) String() string {
	return "reset uses the format: `/points reset`"
}

// Execute implements the Command interface for resetCmd
// It resets all the entries to have 0 points but maintains
// them in the board.Leaderboard.
func (cmd *resetCmd) Execute(ctx context.Context, standings board.Leaderboard) (response Message, err error) {
	if standings == nil {
		return nil, errors.New("could not find leaderboard")
	}
	resetTx := wrapTransacter(standings, standings.Reset)
	if err := executeTransaction(ctx, resetTx); err != nil {
		return nil, err
	}
	return message(resetSuccessText), nil
}

func wrapTransacter(standings board.Leaderboard, f transaction) transaction {
	if ts, ok := standings.(board.Transacter); ok {
		return func(tc context.Context) error {
			return ts.Transact(tc, f)
		}
	}
	return f
}

func executeTransaction(ctx context.Context, tx transaction) error {
	if err := tx(ctx); err != nil {
		log.Errorf(ctx, "failed executing transaction: %+s", err)
		return errors.New(errorText)
	}
	return nil
}
