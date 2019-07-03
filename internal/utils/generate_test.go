package utils_test

import (
	"testing"

	"github.com/tags-drive/core/internal/utils"
)

func BenchmarkGenerate(b *testing.B) {
	const maxSize = 16

	for i := 0; i < b.N; i++ {
		_ = utils.GenerateRandomString(maxSize)
	}
}
