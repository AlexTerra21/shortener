package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

func main() {
	client := resty.New()

	endpoint := "http://localhost:8080/"
	// приглашение в консоли
	fmt.Println("Введите длинный URL")
	// открываем потоковое чтение из консоли
	reader := bufio.NewReader(os.Stdin)
	// читаем строку из консоли
	long, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	long = strings.TrimSuffix(long, "\n")
	// заполняем контейнер данными

	resp, err := client.R().
		SetHeader("Content-Type", "application/text").
		SetBody(long).
		Post(endpoint)

	if err != nil {
		panic(err)
	}
	fmt.Println("Статус-код ", resp.StatusCode())
	fmt.Println(string(resp.Body()))
}
