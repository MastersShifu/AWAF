package utils

import (
	"AuthMicroService/internal/config"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/smtp"
)

type SMTPClient struct {
	client   *smtp.Client
	email    string
	password string
	host     string
	port     string
}

func ConnectToSMTP(cfg *config.Config) (*SMTPClient, error) {
	conn, err := net.Dial("tcp", cfg.SMTP.Host+":"+cfg.SMTP.Port)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к серверу SMTP: %v", err)
	}

	client, err := smtp.NewClient(conn, cfg.SMTP.Host)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании SMTP-клиента: %v", err)
	}

	if err := client.StartTLS(&tls.Config{ServerName: cfg.SMTP.Host}); err != nil {
		return nil, fmt.Errorf("ошибка при запуске STARTTLS: %v", err)
	}

	auth := smtp.PlainAuth("", cfg.SMTP.Email, cfg.SMTP.Password, cfg.SMTP.Host)
	if err := client.Auth(auth); err != nil {
		return nil, fmt.Errorf("ошибка при аутентификации на SMTP-сервере: %v", err)
	}

	return &SMTPClient{
		client:   client,
		email:    cfg.SMTP.Email,
		password: cfg.SMTP.Password,
		host:     cfg.SMTP.Host,
		port:     cfg.SMTP.Port,
	}, nil
}

func (s *SMTPClient) Mail(from string) error {
	return s.client.Mail(from)
}

func (s *SMTPClient) Rcpt(to string) error {
	return s.client.Rcpt(to)
}

func (s *SMTPClient) Data() (io.WriteCloser, error) {
	return s.client.Data()
}

func (s *SMTPClient) SendEmail(to string, code string) error {
	if err := s.Mail(s.email); err != nil {
		return fmt.Errorf("ошибка при указании отправителя: %v", err)
	}

	if err := s.Rcpt(to); err != nil {
		return fmt.Errorf("ошибка при указании получателя: %v", err)
	}

	wc, err := s.Data()
	if err != nil {
		return fmt.Errorf("ошибка при подготовке письма: %v", err)
	}

	subject := "Код авторизации"
	body := code

	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	_, err = wc.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("ошибка при записи данных письма: %v", err)
	}

	err = wc.Close()
	if err != nil {
		return fmt.Errorf("ошибка при закрытии writer: %v", err)
	}

	return nil
}
