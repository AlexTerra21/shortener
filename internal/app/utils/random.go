package utils

import (
	"math/rand"
	"time"
)

var (
	// Переменная для хранения генератора случайных чисел
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	// Массив символов латинского алфавита для генерации строк
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

// Получение случайного целого числа от 1 до 1000
func RandInt() int {
	return rng.Intn(1000)
}
