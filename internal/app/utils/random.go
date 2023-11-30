package utils

import (
	"math/rand"
	"time"
)

// Переменная для хранения генератора случайных чисел
var rng *rand.Rand

// Инициализация генератора случайных чисел
func RandInit() {
	source := rand.NewSource(time.Now().UnixNano())
	rng = rand.New(source)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Генерирует строку из случайных символов латиницы
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}
