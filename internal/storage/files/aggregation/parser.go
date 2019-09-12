package aggregation

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var ErrBadSyntax = errors.New("syntax of a logical expression is incorrect")

// LogicalExpr is a parsed logical expression
type LogicalExpr string

// ParseLogicalExpr returns expression in reverse Polish notation
//
// Valid symbols: digits, !, &, |, (, )
// Examples:
//   - input: "66&!8|7" output: "66 8 ! & 7 |"
//   - input: "(!7|6)&(6|9)" output: "7 ! 6 | 6 9 | &"
//
func ParseLogicalExpr(expr string) (result LogicalExpr, err error) {
	// Just in case
	defer func() {
		if r := recover(); r != nil {
			result = LogicalExpr("")
			err = ErrBadSyntax
		}
	}()

	expr = removeAllSpaces(expr)
	if expr == "" {
		return "", nil
	}

	if !isCorrectExpression([]rune(expr)) {
		return "", ErrBadSyntax
	}

	var (
		builder   strings.Builder
		operators logicalStack
		lastDigit = false
	)

	for _, c := range expr {
		switch {
		case isDigit(c):
			if !lastDigit {
				builder.WriteByte(' ')
			}
			builder.WriteRune(c)
			lastDigit = true
		case c == '&' || c == '|':
			lastDigit = false
			for operators.len > 0 &&
				(operators.top() == '!' || isGreaterPriority(c, operators.top())) {
				builder.WriteByte(' ')
				builder.WriteRune(operators.pop())
			}
			operators.push(c)
		case c == '!':
			lastDigit = false
			operators.push(c)
		case c == '(':
			lastDigit = false
			operators.push(c)
		case c == ')':
			lastDigit = false
			for {
				if operators.len == 0 {
					return "", ErrBadSyntax
				}
				if operators.top() == '(' {
					// Remove '('
					operators.pop()
					break
				}
				builder.WriteByte(' ')
				builder.WriteRune(operators.pop())
			}
		default:
			return "", ErrBadSyntax
		}
	}

	for operators.len > 0 {
		builder.WriteByte(' ')
		builder.WriteRune(operators.pop())
	}

	res := builder.String()

	// Trim first and last spaces
	if res[0] == ' ' {
		res = res[1:]
	}
	if res[len(res)-1] == ' ' {
		res = res[:len(res)-1]
	}

	return LogicalExpr(res), nil
}

func removeAllSpaces(s string) string {
	var b bytes.Buffer
	b.Grow(len(s))

	for _, r := range s {
		if r != ' ' {
			b.WriteRune(r)
		}
	}

	return b.String()
}

func isGreaterPriority(a, b rune) bool {
	switch a {
	case '&':
		return b == '&' || b == '|'
	case '|':
		return b == '&'
	default:
		return false
	}
}

func isCorrectExpression(expr []rune) bool {
	if len(expr) == 0 || !isValidSymbol(expr[0]) {
		return false
	}

	if expr[0] == '&' || expr[len(expr)-1] == '&' ||
		expr[0] == '|' || expr[len(expr)-1] == '|' {
		return false
	}

	invalidPares := []struct {
		f rune
		s rune
	}{
		{'(', ')'},
		{'(', '&'},
		{'(', '|'},
		{')', '('},
		{'&', '&'},
		{'&', '|'},
		{'|', '&'},
		{'|', '|'},
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

	invalidRegexpes := []string{
		`\d!\d`,
		`\d\(`,
	}
	for _, r := range invalidRegexpes {
		reg := regexp.MustCompile(r)
		if reg.MatchString(string(expr)) {
			return false
		}
	}

	return true
}

func isValidSymbol(c rune) bool {
	return ('0' <= c && c <= '9') || c == '!' || c == '&' || c == '|' || c == '(' || c == ')'
}

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}
