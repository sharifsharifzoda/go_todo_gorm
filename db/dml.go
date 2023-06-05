package db

const (
	SelectEmail = `SELECT email FROM users WHERE email = $1`

	CreateUser = `INSERT INTO users (name, email, password)
					VALUES (($1), ($2), ($3)) RETURNING id;`

	CreateTask = `INSERT INTO tasks (name, description, deadline, user_id)
					VALUES (($1), ($2), ($3), ($4)) RETURNING id;`

	//GetAllTasks = `SELECT t.id, t.name, t.description, t.done, t.is_active, u.name
	//				FROM tasks AS t, users AS u
	//				WHERE t.is_active = true AND t.user_id = $1 AND t.user_id = u.id;`

	GetAllTasks = `SELECT t.id, t.name, t.description, t.done, t.is_active, t.deadline, u.name
					FROM tasks AS t INNER JOIN users AS u
					ON t.user_id = u.id
					WHERE t.is_active = true AND t.user_id = $1`

	//GetTaskById = `SELECT t.id, t.name, t.description, t.done, t.is_active, u.name
	//				FROM tasks AS t, users AS u
	//				WHERE t.id = $1 AND t.user_id = $2 AND t.user_id = u.id;`

	GetTaskById = `SELECT t.id, t.name, t.description, t.done, t.is_active, t.deadline, u.name
					FROM tasks AS t INNER JOIN users AS u
					ON t.user_id = u.id
					WHERE t.id = $1 AND t.user_id = $2`

	UpdateTask = `UPDATE tasks SET name = $1, description = $2, done = $3, is_active = $4, deadline = $5
             		WHERE id = $6 AND user_id = $7;`

	DeleteTask = `UPDATE tasks SET is_active = false
             		WHERE id = $1 AND user_id = $2;`

	GetUser = `SELECT id, name, email, password, is_active
					FROM users WHERE email = $1`
)
