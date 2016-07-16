package olilib_network

import (
	"bytes"
	"fmt"
	"strings"
)

// TODO(ian): Replace this with something else
type Cache struct {
	Cache *map[string]string
}

// ExecuteCommandStringList Is a function that makes me terribly sad, as
// generics here would make a world of difference.
func (ctx *ConnectionCtx) ExecuteCommand(command string, args map[string]string) string {
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
			for k, _ := range args {
				val, ok := (*ctx.Cache.Cache)[k]
				if ok {
					retVals[index] = fmt.Sprintf("%s:%s", k, val)
					index++
				}
			}

			return createResponse(command, retVals[0:index])
		}
	case "SET":
		{
			retVals := make([]string, len(args))

			index := 0
			for k, v := range args {
				(*ctx.Cache.Cache)[k] = v

				retVals[index] = fmt.Sprintf("%s:%s", k, v)
				index++
			}

			return createResponse(command, retVals)
		}
	case "REQUEST":
		{
			return ctx.handleRequest(command, args)
		}
	}

	return "Invalid command sent in.\n"
}

func createResponse(command string, retVals []string) string {
	CommandMap := make(map[string]string)
	CommandMap["GET"] = "GOT "
	CommandMap["SET"] = "SAT "

	var buffer bytes.Buffer
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

func (ctx *ConnectionCtx) handleRequest(command string, args map[string]string) string {
	var requestItem string
	// TODO(ian): Support multiple actions per REQUEST in the future.
	for k, _ := range args {
		requestItem = k
		break
	}

	fmt.Println("->", requestItem)
	switch strings.ToUpper(requestItem) {
	case "BLOOMFILTER":
		{
			return (*ctx.Bloomfilter).ConvertToString()
		}
	}

	return "Invalid command sent in.\n"
}
