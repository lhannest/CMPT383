package main

import(
	"fmt"
	"os"
	"io/ioutil"
	"unicode"
)

func containsUnescapedQuote(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '"' && s[i - 1] != '\\' {
				return true
		}
	}

	return false
}

func surroundedByQuotes(s string) bool {
	if len(s) < 2 {
		return false
	} else {
		if s[0] == '"' && s[len(s) - 1] == '"' {
			return true
		} else {
			return false
		}
	}
}

func hasProperEscapeChars(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' {
			n := i + 1
			if s[n] == 'u' {
				if n + 4 < len(s) {
					return isHexDigit(rune(s[n+1])) && isHexDigit(rune(s[n+2])) && isHexDigit(rune(s[n+3])) && isHexDigit(rune(s[n+4]))
				} else {
					return false
				}
			}

			is_valid_rune :=s[n] == '"' ||
							s[n] == '\\' ||
							s[n] == '/' ||
							s[n] == 'b' ||
							s[n] == 'f' ||
							s[n] == 'n' ||
							s[n] == 'r' ||
							s[n] == 't'

			return is_valid_rune
		}
	}
	return true
}

func isString(s string) bool {
	if len(s) == 1 {
		return false
	}

	if surroundedByQuotes(s) {
		middle := s[1 : len(s)-1]
		if containsUnescapedQuote(middle) {
			return false
		}
		if ! hasProperEscapeChars(s) {
			return false
		}

		return true
	}

	return false
}

func isNumber(s string) bool {
	if s[0] == '-' {
		s = s[1:]
	}

	for i, r := range s {
		if r == 'e' || r == 'E' {
			if i != len(s) - 1 {
				after := s[i+1:]
				before := s[:i]
				if after[0] == '-' || after[0] == '+' {
					after = after[1:]
				}

				if before[0] == '-' {
					before = before[1:]
				}

				if len(before) == 0 {
					return false
				}

				return isFloatNumber(before) && isSimpleNumber(after)
			} else {
				return false
			}
		}
	}

	return isFloatNumber(s)
}

func isFloatNumber(s string) bool {
	if len(s) == 0 {
		return false
	}

	if ! unicode.IsNumber(rune(s[0])) {
		return false
	}

	already_encountered := false

	if s[0] == '0' {
		return s[1] == '.'
	} else {
		for _, r := range s {
			if ! unicode.IsNumber(r) && r != '.' {
				return false
			}

			if r == '.' {
				if already_encountered {
					return false
				} else {
					already_encountered = true
				}
			}
		}
	}

	return true
}

func isSimpleNumber(s string) bool {
	for _, r := range s {
		if ! unicode.IsNumber(r) {
			return false
		}
	}

	return true
}

func isHexDigit(r rune) bool {
	// is_lower_case := r == 'a' || r == 'b' || r == 'c' || r == 'd' || r == 'e' || r == 'f'
	// is_upper_case :=  r == 'A' || r == 'B' || r == 'C' || r == 'D' || r == 'E' || r == 'F'
	// is_digit := r == '0' || r == '1' || r == '2' || r == '3' || r == '4' ||
	// 			r == '5' || r == '6' || r == '7' || r == '8' || r == '9'

	// return is_lower_case || is_upper_case || is_digit
	return unicode.In(r, unicode.Properties["ASCII_Hex_Digit"])
}

type token struct {
	instance, name string
	length int
}

func makeToken(instance, name string) token {
	return token{length: len(instance), instance: instance, name: name}
}

func main() {
	bytes, _ := ioutil.ReadAll(os.Stdin)
	input := string(bytes)
	tokens := []token{}

	for len(input) != 0 {
		if input[0] == '\n' {
			input = input[1:]
		} else if input[0] == ' ' {
			// If the next char is empty space, discard it
			input = input[1:]
		} else {
			for i := 1; i <= len(input); i++ {
				substring := input[0:i]

				if substring == "{" {
					t := makeToken(substring, "OPEN_BRACE")
					tokens = append(tokens, t)
				} else if substring == "}" {
					t := makeToken(substring, "CLOSE_BRACE")
					tokens = append(tokens, t)
				} else if substring == "[" {
					t := makeToken(substring, "OPEN_BRACKET")
					tokens = append(tokens, t)
				} else if substring == "]" {
					t := makeToken(substring, "CLOSE_BRACKET")
					tokens = append(tokens, t)
				} else if substring == "," {
					t := makeToken(substring, "COMMA")
					tokens = append(tokens, t)
				} else if substring == "null" {
					t := makeToken(substring, "NULL")
					tokens = append(tokens, t)
				} else if substring == ":" {
					t := makeToken(substring, "COLON")
					tokens = append(tokens, t)
				} else if substring == "true" {
					t := makeToken(substring, "TRUE")
					tokens = append(tokens, t)
				} else if substring == "false" {
					t := makeToken(substring, "FALSE")
					tokens = append(tokens, t)
				} else if isString(substring) {
					t := makeToken(substring, "STRING")
					tokens = append(tokens, t)
				} else if isNumber(substring) {
					t := makeToken(substring, "NUMBER")
					tokens = append(tokens, t)
				}
			}

			
			if len(tokens) == 0 {
				fmt.Println()
				panic("no token found in \"" + input + "\"")
			} else {
				largest := tokens[0]
				for _, t := range tokens {
					if t.length > largest.length {
						largest = t
					}
				}

				display(largest.instance, largest.name)
				// fmt.Println(largest)
				input = input[largest.length:]
				tokens = []token{}
			}
		}
	}
}

func display(token, kind string) {
	fmt.Printf("%-15s %s\n", token, kind)
}