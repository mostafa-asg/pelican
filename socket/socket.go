package socket

import (
	"bufio"
	"net"
	"log"
	"github.com/mostafa-asg/pelican/store"
	"github.com/mostafa-asg/pelican/socket/parser"
	bytesUtil "github.com/mostafa-asg/pelican/util/bytes"
)

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

		if string(command[0:3]) == "PUT" {
			key, value, err := parser.ParsePut(command)
			if err != nil {
				log.Println("Error", err)
				break
			}
			kvStore.Put(key, value)

		} else if string(command[0:3]) == "GET" {
			key, err := parser.ParseGet(command)
			if err != nil {
				log.Println("Error", err)
				break
			}
			value, found := kvStore.Get(key)
			if found {
				con.Write(bytesUtil.ToBytes(uint32(len(value.([]byte)))))
				con.Write(value.([]byte))
			} else {
				con.Write(make([]byte, 4))
			}
		} else if string(command[0:3]) == "DEL" {
			key, err := parser.ParseDel(command)
			if err != nil {
				log.Println("Error", err)
				break
			}
			kvStore.Del(key)
		} else {
			log.Println("Unknown command")
			break
		}
	}

	con.Close()
}