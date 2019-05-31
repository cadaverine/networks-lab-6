## Компьютерные сети. Лабораторная работа VI

Реализовать приложение для пересылки указанного пользователем локального текстового файла на указанный почтовый адрес с использованием любого открытого smtp-сервера.

<hr>

Запуск клиента:

```
go run main.go

Specify the sender email    (enter to default):
Specify the recipient email (enter to default):
Enter password:

Enter subject  (enter to default):
Enter message  (enter to default):
Enter filepath (enter to default):

Result message:

----------------------
Subject: test subject
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary=---BOUNDARY---BOUNDARY---BOUNDARY---BOUNDARY---BOUNDARY---

-----BOUNDARY---BOUNDARY---BOUNDARY---BOUNDARY---BOUNDARY---
Content-Type: text/plain; charset=utf-8

test message


-----BOUNDARY---BOUNDARY---BOUNDARY---BOUNDARY---BOUNDARY---
Content-Type: text/plain; charset=utf-8
Content-Disposition: attachment; filename="README.md"

## Компьютерные сети. Лабораторная работа VI

Реализовать приложение для пересылки указанного пользователем локального текстового файла на указанный почтовый адрес с использованием любого открытого smtp-сервера.

<hr>

-----BOUNDARY---BOUNDARY---BOUNDARY---BOUNDARY---BOUNDARY---

----------------------

Message was successfully sent.
```