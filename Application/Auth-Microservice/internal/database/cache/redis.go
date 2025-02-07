package cache

import (
	"AuthMicroService/internal/config"
	"AuthMicroService/internal/models"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var ctx = context.Background()

type RDClient struct {
	client *redis.Client
}

func ConnectToRedis(cfg *config.Config) (*RDClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
		return nil, err
	}

	log.Println("Connected to Redis:", pong)
	return &RDClient{client: client}, err
}

func (r *RDClient) AddTokensToBlackList(id string, tokens map[string]interface{}) error {
	if err := r.client.HSet(ctx, id, tokens).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RDClient) CheckTokensBlackList(id string, tokens map[string]interface{}) error {
	tokensInfo, err := r.client.HGetAll(ctx, id).Result()
	if err != nil {
		return err
	}

	if tokensInfo["jwt"] == tokens["jwt"] || tokensInfo["rt"] == tokens["rt"] {
		return errors.New("tokens in black list")
	}
	return nil
}

func (r *RDClient) AddCode(email string, code string, credentials models.Credentials) error {
	// Сериализуем структуру в JSON
	credentialsJSON, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	// Используем email как ключ, поле "code" и само значение code
	if err := r.client.HSet(ctx, email, "code", code).Err(); err != nil {
		return err
	}

	// Сохраняем сериализованные credentials
	if err := r.client.HSet(ctx, email, "credentials", credentialsJSON).Err(); err != nil {
		return err
	}

	if err := r.client.Expire(ctx, email, 5*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RDClient) CheckCode(email string, code string) (*models.Credentials, error) {
	// Получаем сохранённый код
	storedCode, err := r.client.HGet(ctx, email, "code").Result()
	if err != nil {
		return nil, err
	}

	// Проверка соответствия кода
	if storedCode != code {
		return nil, errors.New("invalid code")
	}

	// Получаем сохранённые credentials
	credentialsJSON, err := r.client.HGet(ctx, email, "credentials").Result()
	if err != nil {
		return nil, err
	}

	// Десериализуем JSON обратно в структуру
	var credentials models.Credentials
	if err := json.Unmarshal([]byte(credentialsJSON), &credentials); err != nil {
		return nil, err
	}

	// Удаляем запись после успешной проверки
	if err := r.client.Del(ctx, email).Err(); err != nil {
		return nil, err
	}

	return &credentials, nil
}

// UpdateCode обновляет код для указанной почты
func (r *RDClient) UpdateCode(email string, newCode string) error {
	// Проверяем, существует ли ключ в Redis
	exists, err := r.client.Exists(ctx, email).Result()
	if err != nil {
		return err
	}

	if exists == 0 {
		return errors.New("email not found") // Возвращаем ошибку, если ключ не существует
	}

	// Обновляем код
	if err := r.client.HSet(ctx, email, "code", newCode).Err(); err != nil {
		return err
	}

	// Устанавливаем время жизни ключа в 5 минут (если нужно)
	if err := r.client.Expire(ctx, email, 5*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}
