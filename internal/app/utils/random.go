package utils

import (
	"math/rand"
	"time"
)

var (
	// Переменная для хранения генератора случайных чисел
	rng     = rand.New(rand.NewSource(time.Now().UnixNano()))
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// Генерирует строку из случайных символов латиницы
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func RandInt() int {
	return rng.Intn(1000)
}
