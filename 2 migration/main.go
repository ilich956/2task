package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	connStr := "postgres://postgres:bayipket@localhost/adv_database?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to the database")

	defer db.Close()

	dsn := "user=postgres dbname=adv_database sslmode=disable"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	gormDB.AutoMigrate(&User{})

	newUser := User{Name: "Madiyar", Email: "madik@gmail.com"}
	gormDB.Create(&newUser)

	var user User
	gormDB.First(&user, 1)
	log.Println(user)

	gormDB.Model(&user).Update("Name", "Alizhan")

	gormDB.Delete(&user, 1)
}
