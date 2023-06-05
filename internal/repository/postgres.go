package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"todo_gorm/configs"
)

func GetDBConnection(cfg configs.DatabaseConnConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Dushanbe",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName)
	//fmt.Println(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error connecting database. Error is %v", err.Error())
		panic(err.Error())
	}

	log.Printf("Connection success host:%s port:%s", cfg.Host, cfg.Port)

	return db
}

func Close(db *gorm.DB) {
	conn, err := db.DB()
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = conn.Close()
}
