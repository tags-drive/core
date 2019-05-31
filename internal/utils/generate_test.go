package utils_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/tags-drive/core/internal/utils"
)

func BenchmarkGenerate(b *testing.B) {
	const maxSize = 16

	for i := 0; i < b.N; i++ {
		fmt.Fprint(ioutil.Discard, utils.GenerateRandomString(maxSize))
	}
}
