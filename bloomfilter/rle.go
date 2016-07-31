package olilib

import (
	"strconv"
	"strings"
	"fmt"
)

// Encode handles encoding the bloom filter.
// An example string before encoding: "AAAABBCCCDZZZRTTT"
// An example string after encoding "A4B2C3D1Z3R1T3".
// Please note: This could essentially return an inoptimal compression, as if
// we have a string with alternating single values (i.e., "01010101010101")
// this RLE encoding will result in a string twice as long.
// There is an optimization where we only count runs and not single characters,
// but for the time being, I don't think it's necessary.
func Encode(inputString string) string {
	if len(inputString) == 0 {
		return ""
	}

	var currentChar byte
	var count int = 0
	var output string

	// Iterate through each character in the string. Keep the current
	// character that we're tracking and the current count. If we reach a
	// character which is not the currently tracked character, we reset
	// the currently tracked character and the count. Runs in O(n).
	for i := range inputString {
		if i == 0 {
			currentChar = inputString[i]
			count++
			continue
		}

		if inputString[i] == currentChar {
			count++
		} else {
			output = writeOutput(
				output,
				currentChar,
				count,
			)

			currentChar = inputString[i]
			count = 1
		}
	}

	// Write it one more time because we hit bounds on ranging inputString.
	output = writeOutput(
		output,
		currentChar,
		count,
	)


	return output
}

// Decode essentially works opposite of Encode and turns and encoded string
// into a normal, usable string.
func Decode(encodedString string) string {
	if len(encodedString) == 0 {
		return ""
	}

	var output string
	var increment bool

	var i int = 0
	encodedStringLength := len(encodedString) - 1
	for i = 0; i <= encodedStringLength; {
		if i + 1 <= encodedStringLength {
			repeatCount, err := strconv.Atoi(string(encodedString[i + 1]))
			if err != nil {
				repeatCount = 1
				increment = true
			}
			output = writeRepeat(output, encodedString[i], repeatCount)
		} else {
			output = writeRepeat(output, encodedString[i], 1)
		}

		if increment || i + 2 > encodedStringLength + 1 {
			i++
		} else {
			i += 2
		}

		increment = false
	}

	return output
}

// writeOutput is just a simple helper method for Encode which sprintfs the
// output string
func writeOutput(outputString string, char byte, count int) string {
	var retVal string

	if count == 1 {
		return fmt.Sprintf(
			"%s%s",
			outputString,
			string(char),
		)
	} else if count > 9 {
		number := "9"
		retVal = outputString
		for x := 0; x <= count / 9; x++ {
			if x == count / 9  && count % 9 != 0{
				number = strconv.Itoa(count % 9)
			} else if x == count / 9 {
				continue
			}
			retVal = fmt.Sprintf(
				"%s%s%s",
				retVal,
				string(char),
				number,
			)
		}
	} else {
		retVal = fmt.Sprintf(
			"%s%s%d",
			outputString,
			string(char),
			count,
		)
	}


	return retVal
}

// writeRepeat handles writing repeating characters intelligently
func writeRepeat(output string, char byte, repeat int) string {
	var retVal string

	// If the next character is an integer, we can encode it.
	if repeat > 1 {
		retVal = fmt.Sprintf(
			"%s%s",
			output,
			strings.Repeat(
				string(char),
				repeat,
			),
		)
	} else {
		retVal = fmt.Sprintf(
			"%s%s",
			output,
			string(char),
		)
	}

	return retVal
}
