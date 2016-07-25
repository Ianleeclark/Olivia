package incomingNetwork

import (
	"bytes"
	"fmt"
	"github.com/GrappigPanda/Olivia/lib/queryLanguage"
	"strings"
	"github.com/GrappigPanda/Olivia/lib/chord"
)

// ExecuteCommand Is a function that makes me terribly sad, as
// generics here would make a world of difference.
func (ctx *ConnectionCtx) ExecuteCommand(requestData queryLanguage.CommandData) string {
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
				val, ok := (*ctx.Cache.Cache)[k]
				if ok {
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
				(*ctx.Cache.Cache)[k] = v

				retVals[index] = fmt.Sprintf("%s:%s", k, v)
				index++
				(*ctx.Bloomfilter).AddKey([]byte(k))
			}

			return createResponse(command, retVals, requestData.Hash)
		}
	case "REQUEST":
		{
			return ctx.handleRequest(requestData)
		}
	}

	return "Invalid command sent in.\n"
}

func createResponse(command string, retVals []string, hash string) string {
	CommandMap := make(map[string]string)
	CommandMap["GET"] = "GOT "
	CommandMap["SET"] = "SAT "

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

func (ctx *ConnectionCtx) handleRequest(requestData queryLanguage.CommandData) string {
	var requestItem string
	// TODO(ian): Support multiple actions per REQUEST in the future.
	for k := range requestData.Args {
		requestItem = k
		break
	}

	fmt.Println("->", requestItem)
	switch strings.ToUpper(requestItem) {
	case "BLOOMFILTER":
		{
			return (*ctx.Bloomfilter).ConvertToString()
		}
	case "CONNECT":
		{
			peer := chord.NewPeer(requestData.Conn, ctx.MessageBus)
			(*peer).GetBloomFilter()
		}
	}

	return "Invalid command sent in.\n"
}
