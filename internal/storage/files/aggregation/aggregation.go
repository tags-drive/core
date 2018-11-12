package aggregation

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func has(tags []int, tag int) bool {
	for i := range tags {
		if tags[i] == tag {
			return true
		}
	}

	return false
}

// IsGoodFile runs expression for file tags
//
// expr is a logical expression in reverse Polish notation.
//
func IsGoodFile(expr string, fileTags []int) (res bool) {
	// Just in case
	defer func() {
		if r := recover(); r != nil {
			res = false
		}
	}()

	if expr == "" {
		return true
	}

	var (
		steps stack
		s     string
		err   error
	)

	buf := bytes.NewBuffer([]byte(expr))
	for {
		_, err = fmt.Fscan(buf, &s)
		if err == io.EOF {
			break
		}

		switch s {
		case "!":
			b := steps.pop()
			steps.push(!b)
		case "&":
			a := steps.pop()
			b := steps.pop()
			steps.push(a && b)
		case "|":
			a := steps.pop()
			b := steps.pop()
			steps.push(a || b)
		default:
			// We can skip error because expr must be correct
			id, _ := strconv.Atoi(s)
			steps.push(has(fileTags, id))
		}
	}

	if steps.len != 1 {
		return false
	}

	return steps.pop()
}
