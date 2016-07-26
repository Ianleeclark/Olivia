package queryLanguage

import (
	"fmt"
	"github.com/GrappigPanda/Olivia/lru_cache"
	"github.com/GrappigPanda/Olivia/network/message_handler"
	"net"
	"strings"
)

type Parser struct {
	LRUCache     *olilib_lru.LRUCacheString
	MessageStore *message_handler.MessageHandler
}

// CommandData is a struct representing the command sent in.
type CommandData struct {
	Hash    string
	Command string
	Args    map[string]string
	Conn    *net.Conn
}

// NewParser handles creating a new parser (mostly just initializing a new LRU
// cache).
func NewParser(mh *message_handler.MessageHandler) *Parser {
	return &Parser{
		LRUCache: olilib_lru.NewString(25),
	}
}

// Parse handles parsing the grammer into a `CommandData` struct to be later
// processed.
func (p *Parser) Parse(commandString string, conn *net.Conn) (*CommandData, error) {
	splitCommand := strings.SplitN(commandString, " ", 2)
	if len(splitCommand) == 1 {
		return &CommandData{}, fmt.Errorf("%v is an Invalid command.", commandString)
	}

	var hash string
	hashCommand := strings.SplitN(splitCommand[0], ":", 2)
	if len(hashCommand) == 1 {
		hash = ""
	}

	command := splitCommand[0]
	args := make(map[string]string)

	args = parseArgs(strings.Split(splitCommand[1], ","))

	return &CommandData{
		hash,
		command,
		args,
		conn,
	}, nil
}

// parseArgs handles filtering commands based on the command grammer.
// Essentially seperates commands delimited by colons and commands not.
func parseArgs(args []string) map[string]string {
	outMap := make(map[string]string)

	for arg := range args {
		if strings.Contains(args[arg], ":") {
			subCommand := strings.Split(args[arg], ":")
			setKeyValue(&outMap, subCommand[0], subCommand[1])
		} else {
			setKeyValue(&outMap, args[arg], "")
		}
	}

	return outMap
}

// setKeyValue sets a key-value  to a capitalized(key) = value
func setKeyValue(dict *map[string]string, key string, value string) {
	key = strings.Replace(key, "\n", "", -1)
	value = strings.Replace(value, "\n", "", -1)

	(*dict)[strings.ToUpper(key)] = value
}
