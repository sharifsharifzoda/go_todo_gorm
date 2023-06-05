package repository

import (
	"gorm.io/gorm"
	"log"
	"todo_gorm/model"
)

type TodoTaskPostgres struct {
	db *gorm.DB
}

func NewTodoTaskPostgres(db *gorm.DB) *TodoTaskPostgres {
	return &TodoTaskPostgres{db: db}
}

func (r *TodoTaskPostgres) CreateTask(userId int, task *model.Task) (int, error) {
	//var id int
	//row := r.db.QueryRow(db.CreateTask, task.Name, task.Description, task.Deadline, userId)
	//err := row.Scan(&id)
	//if err != nil {
	//	return 0, err
	//}
	//return id, nil

	task.UserId = userId
	tx := r.db.Select("Name", "Description", "Deadline", "UserId").Create(&task)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return task.Id, nil
}

func (r *TodoTaskPostgres) GetAll(userId int) (model.Tasks, error) {
	var tasks model.Tasks
	rows, err := r.db.Model(&model.Task{}).Select([]string{"tasks.id", "tasks.name", "tasks.description",
		"tasks.done", "tasks.is_active", "tasks.deadline", "users.name"}).
		Joins("inner join users on tasks.user_id = users.id").
		Where("tasks.is_active = ? AND tasks.user_id = ?", true, userId).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t model.Task
		err := rows.Scan(&t.Id, &t.Name, &t.Description, &t.Done, &t.IsActive, &t.Deadline, &t.Username)
		if err != nil {
			log.Println("could not to scan from row")
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *TodoTaskPostgres) GetTaskById(userId int, id int) (model.Task, error) {
	var t model.Task
	row := r.db.Model(&model.Task{}).Select([]string{"tasks.id", "tasks.name", "tasks.description", "tasks.done",
		"tasks.is_active", "tasks.deadline", "users.name"}).
		Joins("inner join users on tasks.user_id = users.id").
		Where("tasks.id = ? AND tasks.user_id = ?", id, userId).Row()

	err := row.Scan(&t.Id, &t.Name, &t.Description, &t.Done, &t.IsActive, &t.Deadline, &t.Username)
	if err != nil {
		return model.Task{}, err
	}

	return t, nil
}

func (r *TodoTaskPostgres) UpdateTask(userId int, id int, task *model.Task) error {
	//_, err := r.db.Exec(db.UpdateTask, task.Name, task.Description, task.Done, task.IsActive, task.Deadline, id, userId)
	//return err
	task.Id = id
	tx := r.db.Where("id = ? and user_id = ?", id, userId).Save(task)

	return tx.Error
}

func (r *TodoTaskPostgres) DeleteTask(userId int, id int) error {
	//_, err := r.db.Exec(db.DeleteTask, id, userId)
	//return err
	var task = model.Task{
		Id:       id,
		IsActive: false,
	}
	tx := r.db.Where("id = ? and user_id = ?", id, userId).Save(&task)

	return tx.Error
}
