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
        args := make(map[string]string)

        args = parseArgs(splitCommand[1])

        response := ExecuteCommand(command, args)

        return response, nil
}

func parseArgs(args ...string) map[string]string{
        outMap := make(map[string]string)

        for arg := range args {
                if strings.Contains(args[arg], ":") {
                        subCommand := strings.Split(args[arg], ":") 
                        outMap[subCommand[0]] = subCommand[1]
                }
        }

        return outMap
}
