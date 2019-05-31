package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	smtpHost         = "smtp.yandex.ru"
	smtpPort         = "465"
	defaultSender    = "networks-lab-6@yandex.ru"
	defaultRecipient = "spygelsky@yandex.ru"
	defaultSubject   = "test subject"
	defaultMessage   = "test message"
)

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func handleError(value interface{}, err error) interface{} {
	checkError(err)
	return value
}

func createMessage(headers map[string]string, bodyStr string) string {
	var headersStr string

	for key, value := range headers {
		headersStr += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	return headersStr + "\r\n" + bodyStr
}

func main() {
	body := ""
	sender := ""
	subject := ""
	password := ""
	recipient := ""

	fmt.Print("Specify the sender email    (enter to default): ")
	fmt.Scanln(&sender)

	fmt.Print("Specify the recipient email (enter to default): ")
	fmt.Scanln(&recipient)

	if sender == "" {
		sender = defaultSender
	}

	if recipient == "" {
		recipient = defaultRecipient
	}

	fmt.Print("Enter password: ")
	password = string(handleError(terminal.ReadPassword(int(syscall.Stdin))).([]byte))
	fmt.Print("\n\n")

	tlsConfig := &tls.Config{
		ServerName:         smtpHost,
		InsecureSkipVerify: true,
	}

	auth := smtp.PlainAuth("", sender, password, smtpHost)

	connection := handleError(tls.Dial("tcp", smtpHost+":"+smtpPort, tlsConfig)).(*tls.Conn)

	client := handleError(smtp.NewClient(connection, smtpHost)).(*smtp.Client)

	client.Auth(auth)
	client.Mail(sender)
	client.Rcpt(recipient)

	fmt.Print("Enter subject               (enter to default): ")
	fmt.Scanln(&subject)
	if subject == "" {
		subject = defaultSubject
	}

	fmt.Print("Enter message               (enter to default): ")
	fmt.Scanln(&body)
	if body == "" {
		body = defaultMessage
	}

	headers := make(map[string]string)

	headers["From"] = sender
	headers["To"] = recipient
	headers["Subject"] = subject

	message := createMessage(headers, body)

	fmt.Println("\nResult message:")
	fmt.Println("\n----------------------")
	fmt.Println(message)
	fmt.Println("----------------------")

	writeCloser := handleError(client.Data()).(io.WriteCloser)
	writeCloser.Write([]byte(message))

	fmt.Println("\nMessage was successfully sent.")

	writeCloser.Close()
	client.Quit()
}
