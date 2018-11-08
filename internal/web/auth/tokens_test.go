package auth

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ShoshinNikita/log"
)

func isEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// a in b
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

	// b in a
	for _, i := range b {
		has := false
		for _, j := range a {
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
	// Create tokens.json file in folder web/auth/configs. This file is used for test tokens.write()
	err := os.Mkdir("configs", 0666)
	if err != nil && !os.IsExist(err) {
		log.Fatalln(err)
		return
	}

	f, err := os.Create("configs/tokens.json")
	if err != nil {
		log.Fatalln(err)
		return
	}
	// We have to close the file to remove it
	f.Close()

	// Every test tokens is equal
	code := m.Run()

	// Remove test file
	err = os.Remove("configs/tokens.json")
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = os.Remove("configs")
	if err != nil {
		log.Fatalln(err)
		return
	}

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

	tt.delete("465")
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

func TestExpire(t *testing.T) {
	testTokens := tokens{mutex: new(sync.RWMutex), tokens: originalTokens()}
	tests := []struct {
		before []tokenStruct
		after  []tokenStruct
	}{
		{
			before: []tokenStruct{
				tokenStruct{Token: "123", Expires: time.Now().AddDate(0, 0, -1)},
				tokenStruct{Token: "456", Expires: time.Now().AddDate(0, -2, 0)},
				tokenStruct{Token: "789", Expires: time.Now().AddDate(0, 0, 1)},
			},
			after: []tokenStruct{
				tokenStruct{Token: "789", Expires: time.Now().AddDate(0, 0, 1)},
			},
		},
		{
			before: []tokenStruct{
				tokenStruct{Token: "123", Expires: time.Now().AddDate(1, 2, -1)},
				tokenStruct{Token: "456", Expires: time.Now().AddDate(0, -2, 0)},
				tokenStruct{Token: "789", Expires: time.Now().AddDate(-3, 0, 1)},
			},
			after: []tokenStruct{
				tokenStruct{Token: "123", Expires: time.Now().AddDate(1, 2, -1)},
			},
		},
	}

	for i, tt := range tests {
		testTokens.tokens = make([]tokenStruct, len(tt.before))
		copy(testTokens.tokens, tt.before)
		testTokens.expire()

		want := toStringSlice(tt.after)
		got := toStringSlice(testTokens.tokens)

		if !isEqual(want, got) {
			t.Errorf("Test #%d Want: %v Got: %v\n", i, want, got)
		}
	}
}
