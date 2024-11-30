package main

import (
	"bufio"
	"fmt"
	"os"
)

type char struct {
	val  byte
	meta bool
}

func main() {
	regex, input := parseInput()
	fmt.Println(match(regex, input))
}

func NewChar(val byte, meta bool) char {
	return char{val, meta}
}

func match(regex, input []byte) bool {
	var beginning, end bool
	if len(regex) > 0 && regex[0] == '^' {
		beginning = true
		regex = regex[1:]
	}
	if len(regex) > 0 && regex[len(regex)-1] == '$' {
		end = true
		regex = regex[:len(regex)-1]
	}

	regexChars := bytesToChars(regex, true)

	return invokeMatching(regexChars, bytesToChars(input, false), metachar(regexChars), beginning, end)
}

func invokeMatching(regex, input []char, meta byte, beginning, end bool) bool {
	switch {
	case len(regex) == 0:
		return !end || len(input) == 0
	case len(input) == 0:
		if meta == '?' || meta == '*' {
			return invokeMatching(regex[2:], input, metachar(regex[2:]), beginning, end)
		}
		return len(regex) == 0
	case !beginning && end && len(regex) < len(input):
		return invokeMatching(regex, input[1:], meta, beginning, end)
	default:
		switch meta {
		case '?':
			if equal(regex[0], input[0]) {
				return invokeMatching(regex[2:], input[1:], metachar(regex[2:]), beginning, end)
			}
			return invokeMatching(regex[2:], input, metachar(regex[2:]), beginning, end)
		case '*':
			if !equal(regex[0], input[0]) {
				return invokeMatching(regex[2:], input, metachar(regex[2:]), beginning, end)
			}
			return invokeMatching(regex, input[1:], meta, beginning, end) ||
				invokeMatching(regex, input[1:], '?', beginning, end)
		case '+':
			if equal(regex[0], input[0]) {
				return invokeMatching(regex, input[1:], '*', beginning, end)
			}
			return invokeMatching(regex, input[1:], meta, beginning, end)
		default:
			if equal(regex[0], input[0]) {
				return invokeMatching(regex[1:], input[1:], metachar(regex[1:]), true, end)
			} else if len(regex) == len(input) || beginning {
				return false
			}
			return invokeMatching(regex, input[1:], meta, beginning, end)
		}
	}
}

func equal(r, c char) bool {
	return (r.meta && r.val == '.') || r.val == c.val
}

func metachar(regex []char) byte {
	if len(regex) > 1 && (regex[1].meta && regex[1].val != '.') {
		return regex[1].val
	}

	return 0
}

func bytesToChars(bytes []byte, escape bool) []char {
	var chars []char

	if escape {
		chars = make([]char, 0)
		for len(bytes) > 0 {
			switch bytes[0] {
			case '?', '*', '+', '.':
				chars = append(chars, NewChar(bytes[0], true))
				bytes = bytes[1:]
			case '\\':
				chars = append(chars, NewChar(bytes[1], false))
				bytes = bytes[2:]
			default:
				chars = append(chars, NewChar(bytes[0], false))
				bytes = bytes[1:]
			}
		}
	} else {
		chars = make([]char, len(bytes))
		for i, b := range bytes {
			chars[i] = NewChar(b, false)
		}
	}

	return chars
}

func parseInput() ([]byte, []byte) {
	reader := bufio.NewReader(os.Stdin)

	regex, _ := reader.ReadBytes('|')
	input, _ := reader.ReadBytes('\n')

	return regex[:len(regex)-1], input[:len(input)-1]
}
