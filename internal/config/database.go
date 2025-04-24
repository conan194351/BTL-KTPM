package config

import (
	"fmt"
	"github.com/conan194351/BTL-KTPM/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
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
		), &gorm.Config{
			Logger: gormlog.Default.LogMode(gormlog.Info),
		},
	)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return
	}
	err = db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
		return
	}
	fmt.Println("Postgres connected")
	DB = db
}

func GetDB() *gorm.DB {
	return DB
}
