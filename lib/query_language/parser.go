package query_language

import (
        "strings"
        "fmt"
        "github.com/GrappigPanda/Olivia/lib/lru_cache"
)

// TODO(ian): Add a pointer to the cache here.
type Parser struct {
        LRUCache *olilib_lru.LRUCacheString
}

type CommandData struct {
        Command string
        args []string
}

func (p *Parser) Parse(commandString string) (string, error) {
        splitCommand := strings.SplitN(commandString, " ", 1)
        if len(splitCommand) == 1 {
                return "", fmt.Errorf("%v is an Invalid command.", commandString)
        }

        command := splitCommand[0]

        if command == "SET" {
                subArgs := strings.Split(splitCommand[1], ":")
        } else {
                subArgs := splitCommand[1]
        }

        args := parseArgs(subArgs)

        response, err := ExecuteCommand(command, )

        return response, err
}

func parseArgs(string ...args) {
        
}
