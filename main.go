package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	smtpHost        = "smtp.yandex.ru"
	smtpPort        = "465"
	defaultSubject  = "test subject"
	defaultMessage  = "test message"
	defaultFilePath = "./attachments/test.txt"
	defaultBoundary = "---BOUNDARY---BOUNDARY---BOUNDARY---BOUNDARY---BOUNDARY---"
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

func createMessage(subject, message, fileName, boundary string) []byte {
	buffer := new(bytes.Buffer)

	fmt.Fprintf(buffer, "Subject: %s\n", subject)
	fmt.Fprintf(buffer, "MIME-Version: 1.0\n")
	fmt.Fprintf(buffer, "Content-Type: multipart/mixed; boundary=%s\n", boundary)

	fmt.Fprintf(buffer, "\n--%s\n", boundary)
	fmt.Fprintf(buffer, "Content-Type: text/plain; charset=utf-8\n\n")
	fmt.Fprintf(buffer, "%s\n\n", message)

	if fileName != "" {
		addAttachment(buffer, fileName, boundary)
	}

	return buffer.Bytes()
}

func addAttachment(writer io.Writer, fileName, boundary string) {
	fmt.Fprintf(writer, "\n--%s\n", boundary)
	fmt.Fprintf(writer, "Content-Type: text/plain; charset=utf-8\n")

	fileData, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(writer, "Error: could not open file: %v\n", err)
	} else {
		defer fileData.Close()
		fmt.Fprintf(writer, "Content-Disposition: attachment; filename=\"%s\"\n\n", filepath.Base(fileName))
		io.Copy(writer, fileData)
	}

	fmt.Fprintf(writer, "\n--%s\n", boundary)
}

func main() {
	body := ""
	sender := ""
	subject := ""
	filePath := ""
	password := ""
	recipient := ""

	fmt.Print("Specify the sender email    (enter to default): ")
	fmt.Scanln(&sender)

	fmt.Print("Specify the recipient email (enter to default): ")
	fmt.Scanln(&recipient)

	fmt.Print("Enter password: ")
	password = string(handleError(terminal.ReadPassword(int(syscall.Stdin))).([]byte))
	fmt.Print("\n\n")

	fmt.Print("Enter subject  (enter to default): ")
	fmt.Scanln(&subject)
	if subject == "" {
		subject = defaultSubject
	}

	fmt.Print("Enter message  (enter to default): ")
	fmt.Scanln(&body)
	if body == "" {
		body = defaultMessage
	}

	fmt.Print("Enter filepath (enter to default): ")
	fmt.Scanln(&filePath)
	if filePath == "" {
		filePath = defaultFilePath
	}

	messageBytes := createMessage(subject, body, filePath, defaultBoundary)

	fmt.Println("\nResult message:")
	fmt.Println("\n----------------------")
	fmt.Println(string(messageBytes))
	fmt.Println("----------------------")

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

	data := handleError(client.Data()).(io.WriteCloser)
	defer data.Close()

	io.Copy(data, bytes.NewReader(messageBytes))

	fmt.Println("\nMessage was successfully sent.")
	client.Quit()
}
