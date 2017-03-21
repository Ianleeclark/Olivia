package incomingNetwork

import (
	"bytes"
	"fmt"
	"github.com/GrappigPanda/Olivia/parser"
	"log"
	"strconv"
	"strings"
)

// ExecuteCommand Is a function that makes me terribly sad, as
// generics here would make a world of difference.
func (ctx *ConnectionCtx) ExecuteCommand(requestData parser.CommandData) string {
	command := requestData.Command
	args := requestData.Args

	switch strings.ToUpper(command) {
	case "GET":
		{
			// TODO(ian): This should call a function and if err,
			// lookup the err in a lookup table (a file with a lot
			// of error messages and then return that to the Parser
			// which will return to the parser to the command
			// processor.
			retVals := make([]string, len(args))

			index := 0
			for k := range args {
				val, err := ctx.Cache.Get(k)
				if err == nil {
					retVals[index] = fmt.Sprintf("%s:%s", k, val)
					index++
				}
			}

			return createResponse(command, retVals[0:index], requestData.Hash)
		}
	case "SET":
		{
			retVals := make([]string, len(args))

			index := 0
			for k, v := range args {
				ctx.Cache.Set(k, v)

				retVals[index] = fmt.Sprintf("%s:%s", k, v)
				index++
				ctx.Bloomfilter.AddKey([]byte(k))
			}

			return createResponse(command, retVals, requestData.Hash)
		}
	case "SETEX":
		{
			retVals := make([]string, len(args))
			expirations := requestData.Expiration

			if len(args) != len(expirations) {
				return "Invalid command sent in. Unbalanced keys:expirations.\n"
			}

			index := 0
			for k, v := range args {
				expInt, err := strconv.Atoi(expirations[k])
				if err != nil {
					continue
				}

				log.Println(k, v, expInt)
				(*ctx.Cache).SetExpiration(k, v, expInt)

				retVals[index] = fmt.Sprintf("%s:%s:%d", k, v, expInt)
				index++
				// Please note: Expiration keys are not added to the bloom
				// filter, as the bloom filter only tracks the immutable state
				// of Olivia.
			}

			return createResponse(command, retVals, requestData.Hash)

		}
	case "REQUEST":
		{
			return ctx.handleRequest(requestData)
		}
	case "PING":
		{
			return "0:PONG 1\n"
		}
	}

	return "[]Invalid command sent in.\n"
}

func createResponse(command string, retVals []string, hash string) string {
	CommandMap := make(map[string]string)
	CommandMap["GET"] = "GOT "
	CommandMap["SET"] = "SAT "
	CommandMap["SETEX"] = "SATEX "
	CommandMap["REQUEST"] = "FULFILLED "

	var buffer bytes.Buffer
	buffer.WriteString(hash)
	buffer.WriteString(":")
	buffer.WriteString(CommandMap[command])

	for i := range retVals {
		if i == len(retVals)-1 {
			buffer.WriteString(fmt.Sprintf("%s", retVals[i]))
		} else {
			buffer.WriteString(fmt.Sprintf("%s,", retVals[i]))
		}
	}
	buffer.WriteString("\n")

	return buffer.String()
}

func (ctx *ConnectionCtx) handleRequest(requestData parser.CommandData) string {
	var requestItem string
	// TODO(ian): Support multiple actions per REQUEST in the future.
	for k := range requestData.Args {
		requestItem = k
		break
	}

	switch strings.ToUpper(requestItem) {
	case "BLOOMFILTER":
		{
			bfString := ctx.Bloomfilter.Serialize()
			return createResponse(
				requestData.Command,
				[]string{bfString},
				requestData.Hash,
			)
		}
	case "CONNECT":
		{
			ctx.Cache.AddPeer((*requestData.Conn).RemoteAddr().String())
			return createResponse(
				requestData.Command,
				[]string{ctx.Bloomfilter.Serialize()},
				"",
			)
		}
	case "PEERS":
		{
			return ctx.Cache.ListPeers(requestData.Hash)
		}
	case "DISCONNECT":
		{
			return ctx.Cache.DisconnectPeer((*requestData.Conn).RemoteAddr().String())
		}
	}

	return "Invalid command sent in.\n"
}
