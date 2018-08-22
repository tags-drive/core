package auth

import (
	"fmt"
	"testing"
)

func BenchmarkGenerate(b *testing.B) {
	fmt.Println(b.N)
	for i := 0; i < b.N; i++ {
		fmt.Println(generate(DefaultTokenSize))
	}
}
