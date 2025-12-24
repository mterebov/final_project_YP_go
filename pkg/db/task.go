package db

import "database/sql"


type Task struct {
	ID 		string `json:"id"`
	Date 	string `json:"date"`
	Title 	string `json:"title"`
	Comment string `json:"comment"`
	Repeat 	string `json:"repeat"`
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


func Tasks(limit int64) ([]*Task, error) {
    rows, err := DB.Query(
        `SELECT id, date, title, comment, repeat
         FROM scheduler
         ORDER BY date, id
         LIMIT ?`, limit,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    out := make([]*Task, 0, limit)
    for rows.Next() {
        t := new(Task)
        if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
            return nil, err
        }
        out = append(out, t)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return out, nil
}


// func GetTask(id string) (*Task, error) {
//     query := `SELECT * FROM scheduler WHERE id = :id`
// }


// func UpdateTask(id string) error {
    
// }