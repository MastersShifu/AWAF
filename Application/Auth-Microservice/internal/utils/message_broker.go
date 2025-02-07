package utils

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func SendMessage(username string, id string) {
	// Подключение к RabbitMQ
	conn, err := amqp.Dial("amqp://test:test@auth-rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Ошибка подключения: %s", err)
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	// Создание канала
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Ошибка создания канала: %s", err)
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {

		}
	}(ch)

	// Создание очереди
	q, err := ch.QueueDeclare(
		"hello", // имя очереди
		true,    // долговечная очередь
		false,   // не удалять при завершении работы
		false,   // не эксклюзивная очередь
		false,   // не автоудаление
		nil,
	)

	if err != nil {
		log.Fatalf("Ошибка создания очереди: %s", err)
	}

	// Отправка сообщения
	body := "Username: " + username + " " + "id: " + id

	err = ch.Publish(
		"",     // обменник
		q.Name, // имя очереди
		false,  // обязательное
		false,  // немедленное
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %s", err)
	}

	fmt.Printf("Сообщение отправлено: %s\n", body)
}
