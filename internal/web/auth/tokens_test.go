package auth

import (
	"os"
	"sync"
	"testing"
)

func isEqual(a, b []string) bool {
	for _, i := range a {
		has := false
		for _, j := range b {
			if i == j {
				has = true
				break
			}
		}
		if !has {
			return false
		}
	}

	return true
}

func toStringSlice(t []tokenStruct) (s []string) {
	for _, tt := range t {
		s = append(s, tt.Token)
	}
	return
}

func TestMain(m *testing.M) {
	// Create tokens.json file in folder web/auth. This file is used for test tokens.write()
	f, _ := os.Create("tokens.json")
	// We have to close the file to remove it
	f.Close()

	// Every test tokens is equal
	code := m.Run()

	// Remove test file
	os.Remove("tokens.json")

	os.Exit(code)
}

// originalTokens returns []tokenStruct. The function creates new slice every time.
// It was created to not copy originalTokens every time
func originalTokens() []tokenStruct {
	return []tokenStruct{
		tokenStruct{Token: "123"},
		tokenStruct{Token: "465"},
		tokenStruct{Token: "789"},
		tokenStruct{Token: "101"},
	}
}

func TestAdd(t *testing.T) {
	tt := tokens{mutex: new(sync.RWMutex), tokens: originalTokens()}

	tt.add("999")
	answerSlice := []tokenStruct{
		tokenStruct{Token: "123"},
		tokenStruct{Token: "465"},
		tokenStruct{Token: "789"},
		tokenStruct{Token: "101"},
		tokenStruct{Token: "999"},
	}
	want := toStringSlice(answerSlice)
	got := toStringSlice(tt.tokens)
	if !isEqual(want, got) {
		t.Errorf("Wrong add result Want: %v Got: %v", want, got)
	}

	tt.add("15")
	answerSlice = []tokenStruct{
		tokenStruct{Token: "123"},
		tokenStruct{Token: "465"},
		tokenStruct{Token: "789"},
		tokenStruct{Token: "101"},
		tokenStruct{Token: "999"},
		tokenStruct{Token: "15"},
	}
	want = toStringSlice(answerSlice)
	got = toStringSlice(tt.tokens)
	if !isEqual(want, got) {
		t.Errorf("Wrong add result Want: %v Got: %v", want, got)
	}
}

func TestDelete(t *testing.T) {
	tt := tokens{mutex: new(sync.RWMutex), tokens: originalTokens()}

	tt.delete("456")
	answerSlice := []tokenStruct{
		tokenStruct{Token: "123"},
		tokenStruct{Token: "789"},
		tokenStruct{Token: "101"},
	}
	want := toStringSlice(answerSlice)
	got := toStringSlice(tt.tokens)
	if !isEqual(want, got) {
		t.Errorf("Wrong delete result Want: %v Got: %v", want, got)
	}

	tt.delete("123")
	answerSlice = []tokenStruct{
		tokenStruct{Token: "789"},
		tokenStruct{Token: "101"},
	}
	want = toStringSlice(answerSlice)
	got = toStringSlice(tt.tokens)
	if !isEqual(want, got) {
		t.Errorf("Wrong delete result Want: %v Got: %v", want, got)
	}

	tt.delete("789")
	answerSlice = []tokenStruct{
		tokenStruct{Token: "101"},
	}
	want = toStringSlice(answerSlice)
	got = toStringSlice(tt.tokens)
	if !isEqual(want, got) {
		t.Errorf("Wrong delete result Want: %v Got: %v", want, got)
	}

	tt.delete("101")
	answerSlice = []tokenStruct{}
	want = toStringSlice(answerSlice)
	got = toStringSlice(tt.tokens)
	if !isEqual(want, got) {
		t.Errorf("Wrong delete result Want: %v Got: %v", want, got)
	}

	tt.delete("999")
	answerSlice = []tokenStruct{}
	want = toStringSlice(answerSlice)
	got = toStringSlice(tt.tokens)
	if !isEqual(want, got) {
		t.Errorf("Wrong delete result Want: %v Got: %v", want, got)
	}
}

func TestCheck(t *testing.T) {
	tt := tokens{mutex: new(sync.RWMutex), tokens: originalTokens()}

	res := tt.check("15")
	answerBool := false
	if res != answerBool {
		t.Errorf("Wrong check result Want: %v Got: %v", answerBool, res)
	}

	res = tt.check("123")
	answerBool = true
	if res != answerBool {
		t.Errorf("Wrong check result Want: %v Got: %v", answerBool, res)
	}
}
