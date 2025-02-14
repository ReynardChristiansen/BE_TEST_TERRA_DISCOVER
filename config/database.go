package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "github.com/lib/pq"
)

var DB *gorm.DB

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Gagal memuat file .env:", err)
	}
}

func CreateDatabase() {
	LoadEnv()

	dbConnStr := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	sqlDB, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal("Gagal koneksi ke PostgreSQL:", err)
	}
	defer sqlDB.Close()

	var exists bool
	sqlDB.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)", os.Getenv("DB_NAME")).Scan(&exists)

	if !exists {
		_, err = sqlDB.Exec("CREATE DATABASE " + os.Getenv("DB_NAME"))
		if err != nil {
			log.Fatal("Gagal membuat database:", err)
		}
		log.Println("Database `" + os.Getenv("DB_NAME") + "` berhasil dibuat!")
	} else {
		log.Println("Database `" + os.Getenv("DB_NAME") + "` sudah ada!")
	}
}

func ConnectDB() {
	LoadEnv()

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	log.Println("Database `" + os.Getenv("DB_NAME") + "` Connected")
	DB = db
}
