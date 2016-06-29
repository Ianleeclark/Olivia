package query_language

import (
        "strings"
)

// ExecuteCommandStringList Is a function that makes me terribly sad, as
// generics here would make a world of difference.
func ExecuteCommand(command string, args map[string]string) string {
        switch strings.ToUpper(command) {
                case "GET": {
                        // TODO(ian): This should call a function and if err,
                        // lookup the err in a lookup table (a file with a lot
                        // of error messages and then return that to the Parser
                        // which will return to the parser to the command
                        // processor.
                        
                }
        }
}
