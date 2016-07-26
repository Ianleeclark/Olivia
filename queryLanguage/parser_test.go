package queryLanguage

import (
	"github.com/GrappigPanda/Olivia/network/message_handler"
	"testing"
)

var MESSAGEHANDLER = message_handler.NewMessageHandler()

func TestParseFailInvalidCommand(t *testing.T) {
	parser := NewParser(MESSAGEHANDLER)

	_, err := parser.Parse("XYZalsdkj", nil)
	if err == nil {
		t.Fatalf("Command was somehow parsed as correct.")
	}
}

func TestParseStringOfCommas(t *testing.T) {
	args := make(map[string]string)
	args["key1"] = ""
	args["key2"] = ""
	args["key3"] = ""

	expectedReturn := &CommandData{
		"",
		"GET",
		args,
		nil,
	}

	parser := NewParser(MESSAGEHANDLER)

	retval, err := parser.Parse("GET key1,key2,key3", nil)
	if err != nil {
		t.Fatalf("Failed to parse string `GET key1, key2, key3` with error: %v", err)
	}

	if (*retval).Command != (*expectedReturn).Command {
		t.Fatalf("Expected command %v, got %v", (*retval).Command, (*expectedReturn).Command)
	}

	for key := range expectedReturn.Args {
		if (*retval).Args[key] != (*expectedReturn).Args[key] {
			t.Fatalf("Expected %v, got %v", expectedReturn, retval)
		}
	}
}

func TestParseSetKeysWithColon(t *testing.T) {
	args := make(map[string]string)
	args["key1"] = ""
	args["key2"] = ""
	args["key3"] = ""

	expectedReturn := &CommandData{
		"",
		"SET",
		args,
		nil,
	}

	parser := NewParser(MESSAGEHANDLER)

	retval, err := parser.Parse("GET key1,key2,key3", nil)
	if err != nil {
		t.Fatalf("Failed to parse string `GET key1, key2, key3` with error: %v", err)
	}

	for key := range expectedReturn.Args {
		if (*retval).Args[key] != (*expectedReturn).Args[key] {
			t.Fatalf("Expected %v, got %v", expectedReturn, retval)
		}
	}
}

func TestParseCommandWithHash(t *testing.T) {
	args := make(map[string]string)
	args["key1"] = ""
	args["key2"] = ""
	args["key3"] = ""

	expectedReturn := &CommandData{
		"",
		"SET",
		args,
		nil,
	}

	parser := NewParser(MESSAGEHANDLER)

	retval, err := parser.Parse("GET key1,key2,key3", nil)
	if err != nil {
		t.Fatalf("Failed to parse string `GET key1, key2, key3` with error: %v", err)
	}

	for key := range expectedReturn.Args {
		if (*retval).Args[key] != (*expectedReturn).Args[key] {
			t.Fatalf("Expected %v, got %v", expectedReturn, retval)
		}
	}

}
