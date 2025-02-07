package postgresql

import (
	"AuthMicroService/internal/config"
	"AuthMicroService/internal/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBClient struct {
	database *gorm.DB
}

func ConnectDB(cfg *config.Config) (*DBClient, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&models.Users{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &DBClient{database: db}, nil
}

func (db *DBClient) Table(name string) *gorm.DB {
	return db.database.Table(name)
}

func (db *DBClient) AuthenticateUser(name, password string) (*models.Users, error) {
	var user models.Users

	err := db.Table("users").Where("name = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}

func (db *DBClient) CheckUserExist(name string, email string) error {
	var user models.Users

	err := db.Table("users").Where("name = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = db.Table("users").Where("email = ?", email).First(&user).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return errors.New("user exist")
		}
		return err
	}

	return errors.New("user exist")
}

func (db *DBClient) CheckUserEmailExist(email string) error {
	err := db.Table("users").Where("email = ?", email).First(&models.Users{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DBClient) RegisterNewUser(credentials *models.Credentials) error {
	credentials.ID = uuid.NewString()

	if err := db.Table("users").Create(&credentials).Error; err != nil {
		return err
	}

	return nil
}

func (db *DBClient) GetCredentialsVersion(id string, version int) error {
	var credentials models.Credentials

	if err := db.Table("users").Where("id = ?", id).First(&credentials).Error; err != nil {
		return err
	}

	if version != credentials.Version {
		return errors.New("invalid version")
	}

	return nil
}
