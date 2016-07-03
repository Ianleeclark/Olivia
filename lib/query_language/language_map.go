package query_language

import (
        "strings"
        "bytes"
        "fmt"
)

// TODO(ian): Replace this with something else
type Cache struct {
       Cache *map[string]string
}

// ExecuteCommandStringList Is a function that makes me terribly sad, as
// generics here would make a world of difference.
func (c *Cache) ExecuteCommand(command string, args map[string]string) string {
        switch strings.ToUpper(command) {
                case "GET": {
                        // TODO(ian): This should call a function and if err,
                        // lookup the err in a lookup table (a file with a lot
                        // of error messages and then return that to the Parser
                        // which will return to the parser to the command
                        // processor.
                        retVals := make([]string, len(args))

                        index := 0
                        for k, _ := range args {
                                val, ok := (*c.Cache)[k]
                                if ok {
                                        retVals[index] = fmt.Sprintf("%s:%s", k, val)
                                        index++
                                }
                        }

                        return createResponse(command, retVals[0:index])
                }
                case "SET": {
                        retVals := make([]string, len(args))

                        index := 0
                        for k, v := range args {
                                (*c.Cache)[k] = v

                                retVals[index] = fmt.Sprintf("%s:%s", k, v)
                                index++
                        }

                        return createResponse(command, retVals)
                }
        }

        return "Invalid command sent in."
}

func createResponse(command string, retVals []string) string {
        CommandMap := make(map[string]string)
        CommandMap["GET"] = "GOT "
        CommandMap["SET"] = "SAT "

        var buffer bytes.Buffer
        buffer.WriteString(CommandMap[command])

        for i := range retVals {
                if i == len(retVals) - 1 {
                        buffer.WriteString(fmt.Sprintf("%s", retVals[i]))
                } else {
                        buffer.WriteString(fmt.Sprintf("%s,", retVals[i]))
                }
        }

        return buffer.String()
}
