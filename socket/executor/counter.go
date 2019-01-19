package executor

import (
	"net"

	"github.com/mostafa-asg/pelican/socket/parser"
	"github.com/mostafa-asg/pelican/store"
	bytesUtil "github.com/mostafa-asg/pelican/util/bytes"
)

type CounterExecutor struct{}

func (e *CounterExecutor) Accept(result parser.ParseResult) bool {
	return result.Command() == "C"
}

func (e *CounterExecutor) Execute(result parser.ParseResult, store *store.Store, con net.Conn) {
	r := result.(*parser.CounterResult)

	var counter int64
	if r.Typ == "Inc" {
		counter = store.IncCounter(r.Key, int64(r.Change))
	} else {
		counter = store.DecCounter(r.Key, int64(r.Change))
	}

	con.Write(bytesUtil.ToBytes64(uint64(counter)))
}

func NewCounterExecutor() *CounterExecutor {
	return &CounterExecutor{}
}
