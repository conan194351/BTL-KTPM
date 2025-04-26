package config

import (
	"fmt"
	"github.com/conan194351/BTL-KTPM/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Database struct {
	Host string `env:"HOST" `
	Port string `env:"PORT"`
	User string `env:"USER" `
	Pass string `env:"PASS"`
	Name string `env:"NAME" `
}

var DB *gorm.DB

func (d Database) GetDSN() string {
	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Pass, d.Name)
	return dsn
}

func InitDatabase() {
	cnf := GetConfig().Database
	source := cnf.GetDSN()
	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN:                  source,
				PreferSimpleProtocol: true,
			},
		),
	)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}
	err = db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
		return
	}
	err = seedData(db)
	if err != nil {
		log.Fatal("Failed to seed data:", err)
		return
	}

	fmt.Println("Postgres connected")
	DB = db
}

func seedData(db *gorm.DB) error {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		users := []models.User{
			{Name: "Phu Nguyen", Email: "longphunguyen1404@gmail.com", Password: "hashedpassword1", Phone: "0912345678"},
			{Name: "Bob Tran", Email: "bob@example.com", Password: "hashedpassword2", Phone: "0987654321"},
			{Name: "Charlie Le", Email: "charlie@example.com", Password: "hashedpassword3", Phone: "0909090909"},
		}
		if err := db.Create(&users).Error; err != nil {
			return fmt.Errorf("failed to seed users: %w", err)
		}
	}

	// Check and seed products
	db.Model(&models.Product{}).Count(&count)
	if count == 0 {
		products := []models.Product{
			{Name: "MacBook Pro 14", Description: "Apple M3 Pro, 16GB RAM, 512GB SSD", Price: 1999.99, Stock: 10, ImageURL: "https://example.com/macbook.jpg"},
			{Name: "iPhone 15 Pro", Description: "128GB, Titanium Blue", Price: 999.99, Stock: 25, ImageURL: "https://example.com/iphone.jpg"},
			{Name: "AirPods Pro", Description: "Wireless Noise Cancelling Earbuds", Price: 249.99, Stock: 50, ImageURL: "https://example.com/airpods.jpg"},
		}
		if err := db.Create(&products).Error; err != nil {
			return fmt.Errorf("failed to seed products: %w", err)
		}
	}

	return nil
}

func GetDB() *gorm.DB {
	return DB
}
