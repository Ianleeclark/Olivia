package incomingNetwork

import (
	"bytes"
	"fmt"
	"github.com/GrappigPanda/Olivia/dht"
	"github.com/GrappigPanda/Olivia/parser"
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
			responseChannel := make(chan string)
			for k := range args {
				val, ok := (*ctx.Cache.Cache)[k]
				if ok {
					retVals[index] = fmt.Sprintf("%s:%s", k, val)
					index++
				} else {
					for _, peer := range ctx.PeerList.Peers {
						if peer == nil || peer.Status == dht.Timeout || peer.Status == dht.Disconnected {
							continue
						}

						peer.SendRequest(
							fmt.Sprintf("GET %s", k),
							responseChannel,
							ctx.MessageBus,
						)

						value := <-responseChannel
						if value != "" {
							retVals[index] = fmt.Sprintf("%s:%s", k, value)
							index++
						}
					}
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
			bfString := (*ctx.Bloomfilter).ConvertToString()
			return createResponse(
				requestData.Command,
				[]string{bfString},
				requestData.Hash,
			)
		}
	case "CONNECT":
		{
			(*ctx.PeerList).AddPeer((*requestData.Conn).RemoteAddr().String())
			return createResponse(
				requestData.Command,
				[]string{(*ctx.Bloomfilter).ConvertToString()},
				"",
			)
		}
	case "PEERS":
		{
			count := 0
			outString := fmt.Sprintf("%s:FULFILLED ", requestData.Hash)

			for _, peer := range ctx.PeerList.Peers {
				if peer == nil {
					continue
				}

				if count == 0 {
					outString = fmt.Sprintf(
						"%s%s",
						outString,
						peer.IPPort,
					)

				} else {
					outString = fmt.Sprintf(
						"%s,%s",
						outString,
						peer.IPPort,
					)

				}
			}

			for _, peer := range ctx.PeerList.BackupPeers {
				if peer == nil {
					continue
				}

				outString = fmt.Sprintf(
					"%s,%s",
					outString,
					peer.IPPort,
				)
			}

			outString = fmt.Sprintf(
				"%s\n",
				outString,
			)

			return outString
		}
	case "DISCONNECT":
		{
			outString := "Peer not found in peer list."
			for _, peer := range ctx.PeerList.Peers {
				if peer.IPPort != (*requestData.Conn).RemoteAddr().String() {
					continue
				}

				if peer != nil && peer.Status == dht.Connected {
					peer.Disconnect()
					outString = "Peer has been disconnected."
				}

				// TODO(ian): Connect a backup node after one node has forced itself to be evicted.
			}

			return outString
		}
	}

	return "Invalid command sent in.\n"
}
