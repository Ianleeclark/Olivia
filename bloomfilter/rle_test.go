package olilib

import (
	"testing"
)

func TestEncode(t *testing.T) {
	expectedReturn := "A3B2Z5T3"
	retVal := Encode("AAABBZZZZZTTT")

	if expectedReturn != retVal {
		t.Errorf("Expected %v, got %v", expectedReturn, retVal)
	}
}

func TestEncodeIntegers(t *testing.T) {
	expectedReturn := "052234"
	retVal := Encode("00000223333")

	if expectedReturn != retVal {
		t.Errorf("Expected %v, got %v", expectedReturn, retVal)
	}
}
func TestDecode(t *testing.T) {
	expectedReturn := "AAABBZZZZZTTT"
	retVal := Decode("A3B2Z5T3")

	if expectedReturn != retVal {
		t.Errorf("Expected %v, got %v", expectedReturn, retVal)
	}
}

func TestDencodeIntegers(t *testing.T) {
	expectedReturn := "00000223333"
	retVal := Decode("052234")

	if expectedReturn != retVal {
		t.Errorf("Expected %v, got %v", expectedReturn, retVal)
	}
}

func TestEncodeDecode(t *testing.T) {
	expectedReturn := "AAABBZZZZZTTT"
	retVal := Encode(expectedReturn)
	retVal = Decode(retVal)

	if expectedReturn != retVal {
		t.Errorf("Expected %v, got %v", expectedReturn, retVal)
	}
}

func TestEncodeDecodeIntegers(t *testing.T) {
	expectedReturn := "00000223333"
	retVal := Encode(expectedReturn)
	retVal = Decode(retVal)

	if expectedReturn != retVal {
		t.Errorf("Expected %v, got %v", expectedReturn, retVal)
	}
}
