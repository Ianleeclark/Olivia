package parser

import (
	"fmt"
	"net"
	"strings"

	"github.com/GrappigPanda/Olivia/network/message_handler"
)

type Parser struct {
	LRUCache     *olilib_lru.LRUCacheString
	MessageStore *message_handler.MessageHandler
}

// CommandData is a struct representing the command sent in.
type CommandData struct {
	Hash       string
	Command    string
	Args       map[string]string
	Expiration map[string]string
	Conn       *net.Conn
}

// NewParser handles creating a new parser (mostly just initializing a new LRU
// cache).
func NewParser(mh *message_handler.MessageHandler) *Parser {
	return &Parser{
		LRUCache: olilib_lru.NewString(25),
	}
}

// Parse handles parsing the grammar into a `CommandData` struct to be later
// processed.
func (p *Parser) Parse(commandString string, conn *net.Conn) (*CommandData, error) {
	splitCommand := strings.SplitN(commandString, " ", 2)
	if len(splitCommand) == 1 {
		return &CommandData{}, fmt.Errorf("%v is an Invalid command.", commandString)
	}

	var hash string
	var command string
	hashAndCommand := strings.Split(splitCommand[0], ":")
	if len(hashAndCommand) >= 2 {
		hash = hashAndCommand[0]
		command = hashAndCommand[1]

	} else if len(hashAndCommand) == 1 {
		hash = ""
		command = hashAndCommand[0]
	}

	args, expirations := parseArgs(strings.Split(splitCommand[1], ","))

	return &CommandData{
		hash,
		command,
		args,
		expirations,
		conn,
	}, nil
}

// parseArgs handles filtering commands based on the command grammar.
// Essentially separates commands delimited by colons and commands not.
func parseArgs(args []string) (map[string]string, map[string]string) {
	argMap := make(map[string]string)
	expirationMap := make(map[string]string)

	for arg := range args {
		if strings.Contains(args[arg], ":") {
			subCommand := strings.Split(args[arg], ":")
			setKeyValue(&argMap, subCommand[0], subCommand[1])

			if len(subCommand) > 2 {
				setKeyValue(&expirationMap, subCommand[0], subCommand[2])
			}
		} else {
			setKeyValue(&argMap, args[arg], "")
		}
	}

	return argMap, expirationMap
}

// setKeyValue sets a key-value  to a capitalized(key) = value
func setKeyValue(dict *map[string]string, key string, value string) {
	key = strings.Replace(key, "\n", "", -1)
	value = strings.Replace(value, "\n", "", -1)

	(*dict)[key] = value
}
