package executor

import (
	"net"

	"github.com/mostafa-asg/pelican/socket/parser"
	"github.com/mostafa-asg/pelican/store"
)

type Executor interface {
	Accept(parser.ParseResult) bool
	Execute(parser.ParseResult, *store.Store, net.Conn)
}
