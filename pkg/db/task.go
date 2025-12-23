package db

import "database/sql"


type Task struct {
	ID 		string `json:"id"`
	Date 	string `json:"date"`
	Title 	string `json:"title"`
	Comment string `json:"comment"`
	Repeat 	string `json:"repeat"`
	Error 	string `json:"error"`
}


func AddTask(task *Task) (int64, error) {
	var id int64
    query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
    res, err := DB.Exec(query, 
								sql.Named("date", task.Date),
								sql.Named("title", task.Title),
								sql.Named("comment", task.Comment),
								sql.Named("repeat", task.Repeat),
						)
    if err == nil {
        id, err = res.LastInsertId()
    }
    return id, err
}
