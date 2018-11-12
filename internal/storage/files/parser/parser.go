package parser

import (
	"github.com/pkg/errors"
)

var ErrBadSyntax = errors.New("bad syntax")

// Parse returns exprs in reverse Polish notation
//
// Valid symbols: digits, !, &, |, (, )
// Examples:
//   - input: "66&!8|7" output: "66 8 ! & 7 |"
//   - input: "(!7|6)&(6|9)" output: "7 ! 6 | 6 9 | &"
//
func Parse(expr string) (res string, err error) {
	// Just in case
	defer func() {
		if r := recover(); r != nil {
			res = ""
			err = ErrBadSyntax
		}
	}()

	if !isCorrectExpression(expr) {
		return "", ErrBadSyntax
	}

	var operators stack

	lasDigit := false

	for i := 0; i < len(expr); i++ {
		c := expr[i]
		switch {
		case c == '&' || c == '|':
			lasDigit = false
			for operators.len > 0 &&
				(operators.top() == '!' || isGreaterPriority(c, operators.top())) {
				res += " " + string(operators.pop())
			}
			operators.push(c)
		case c == '!':
			lasDigit = false
			operators.push(c)
		case c == '(':
			lasDigit = false
			operators.push(c)
		case c == ')':
			lasDigit = false
			for {
				if operators.len == 0 {
					return "", ErrBadSyntax
				}
				if operators.top() == '(' {
					// Remove '('
					operators.pop()
					break
				}
				res += " " + string(operators.pop())
			}

		case isDigit(c):
			if !lasDigit {
				res += " "
			}
			res += string(c)
			lasDigit = true
		default:
			return "", ErrBadSyntax
		}
	}

	for operators.len > 0 {
		res += " " + string(operators.pop())
	}

	// Trim first and last spaces
	if res[0] == ' ' {
		res = res[1:]
	}
	if res[len(res)-1] == ' ' {
		res = res[:len(res)-1]
	}

	return res, nil
}

func isGreaterPriority(a, b byte) bool {
	switch a {
	case '&':
		return b == '&' || b == '|'
	case '|':
		return b == '&'
	default:
		return false
	}
}

func isCorrectExpression(expr string) bool {
	if len(expr) == 0 || !isValidSymbol(expr[0]) {
		return false
	}

	if expr[0] == '&' || expr[len(expr)-1] == '&' ||
		expr[0] == '|' || expr[len(expr)-1] == '|' {
		return false
	}

	invalidPares := []struct {
		f byte
		s byte
	}{
		{'(', ')'},
		{'(', '&'},
		{'(', '|'},
		{')', '('},
		{'&', '|'},
		{'|', '&'},
		{'!', '&'},
		{'!', '|'},
		{'!', '!'},
	}

	for i := 1; i < len(expr); i++ {
		if !isValidSymbol(expr[i]) {
			return false
		}

		for _, pare := range invalidPares {
			if expr[i-1] == pare.f && expr[i] == pare.s {
				return false
			}
		}
	}

	return true
}

func isValidSymbol(c byte) bool {
	return ('0' <= c && c <= '9') || c == '!' || c == '&' || c == '|' || c == '(' || c == ')'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isOperator(c byte) bool {
	return c == '&' || c == '|' || c == '!' || c == '(' || c == ')'
}
