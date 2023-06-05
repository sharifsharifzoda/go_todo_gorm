package db

import (
	"gorm.io/gorm"
	"log"
	"todo_gorm/model"
)

func Init(db *gorm.DB) {
	//DDls := []string{
	//	CreateUsersTable,
	//	CreateTasksTable,
	//}
	//
	//for i, ddl := range DDls {
	//	tx := db.Exec(ddl)
	//	if tx.Error != nil {
	//		log.Fatalf("failed to create table #%d due to: %s", i+1, tx.Error.Error())
	//	}
	//}

	err := db.AutoMigrate(&model.User{}, &model.Task{})
	if err != nil {
		log.Fatal(err)
	}
}
