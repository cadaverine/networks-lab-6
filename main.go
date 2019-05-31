package main

import (
	"fmt"
)

const (
	defaultSMTPServer = "smtp.yandex.ru:465"
)

func main() {
	var sender, recipient string

	fmt.Print("Specify the sender email: ")
	fmt.Scanln(&sender)

	fmt.Print("Specify the recipient email: ")
	fmt.Scanln(&recipient)

	fmt.Println(sender, recipient)
}
