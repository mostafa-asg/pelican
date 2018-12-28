package socket

import (
	"bufio"
	"log"
	"net"

	"github.com/mostafa-asg/pelican/socket/executor"
	"github.com/mostafa-asg/pelican/socket/parser"
	"github.com/mostafa-asg/pelican/store"
	bytesUtil "github.com/mostafa-asg/pelican/util/bytes"
)

var parsers []parser.Parser
var executors []executor.Executor

func init() {
	parsers = make([]parser.Parser, 0)
	parsers = append(parsers, parser.NewPutParser())
	parsers = append(parsers, parser.NewGetParser())
	parsers = append(parsers, parser.NewDelParser())

	executors = make([]executor.Executor, 0)
	executors = append(executors, executor.NewPutExecutor())
	executors = append(executors, executor.NewGetExecutor())
	executors = append(executors, executor.NewDelExecutor())
}

func HandleConnection(con net.Conn, kvStore *store.Store) {
	reader := bufio.NewReader(con)

	for {
		bytes, err := bytesUtil.ReadBytes(4, reader)
		if err != nil {
			log.Println("Error", err)
			break
		}

		commandSize := bytesUtil.ToUint32(bytes)
		if commandSize == 0 {
			break
		}

		command, err := bytesUtil.ReadBytes(int(commandSize), reader)
		if err != nil {
			log.Println("Error", err)
			break
		}

		commandFound := false
		for _, p := range parsers {
			if p.Accept(command) {
				commandFound = true

				pr, err := p.Parse(command)
				if err != nil {
					log.Println("Error", err)
					break
				}

				for _, e := range executors {
					if e.Accept(pr) {
						e.Execute(pr, kvStore, con)
						break
					}
				}

				break
			}
		}

		if !commandFound {
			log.Println("Error: Command not found.")
			break
		}
	}

	con.Close()
}
