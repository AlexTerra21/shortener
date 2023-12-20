package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandom_RandSeq(t *testing.T) {
	RandInit()
	randomStrings := make(map[string]int)

	charLimit := 8
	stringLimit := 100
	// Проверка на несовпадение выборки из 100 сгенерированных строк.
	for i := 0; i < stringLimit; i++ {
		randomStrings[RandSeq(charLimit)]++
	}
	for _, count := range randomStrings {
		require.Equal(t, count, 1)
	}

}
