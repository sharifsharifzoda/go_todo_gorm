package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"todo_gorm"
	"todo_gorm/configs"
	"todo_gorm/db"
	"todo_gorm/internal/handler"
	"todo_gorm/internal/repository"
	"todo_gorm/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	//reading from yaml
	if err := InitConfigs(); err != nil {
		log.Fatalf("error while reading config file. error is %v", err.Error())
	}

	var cfg configs.DatabaseConnConfig

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Couldn't unmarshal the config into struct. error is %v", err.Error())
	}
	cfg.Password = os.Getenv("DB_PASSWORD")

	conn := repository.GetDBConnection(cfg)

	db.Init(conn)

	//---------- Dependency injection-----------
	newRepository := repository.NewRepository(conn)
	newService := service.NewService(newRepository)
	newHandler := handler.NewHandler(newService.Auth, newService.Todo)
	//--------------------------------------------

	server := new(todo_sql.Server)
	go func() {
		if err := server.Run(os.Getenv("PORT"), newHandler.InitRoutes()); err != nil {
			log.Fatalf("error while running http.server. Error is %s", err.Error())
		}
	}()
	fmt.Printf("Server is listening to port: %s\n", os.Getenv("PORT"))

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	repository.Close(conn)

	fmt.Println("server is shutting down")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("error while shutting server down. Error: %s", err.Error())
	}
}

func InitConfigs() error {
	viper.AddConfigPath("configs") //адрес директории
	viper.SetConfigName("config")  //имя файла
	viper.SetConfigType("yml")
	return viper.ReadInConfig() //считывает config и сохраняет данные во внутренний объект viper
}
