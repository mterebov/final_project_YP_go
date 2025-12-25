package db

import (
	"database/sql"
)


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


func Tasks(limit int64, search string) ([]*Task, error) {
    var queryString string
    var args []any
    switch search {
        case "":
            queryString = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date, id LIMIT ?`
            args = append(args, limit)
        default:
            if date, err := searchToDate(search); err == nil {
                queryString = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
                args = append(args, date, limit)
            } else {
                search = "%" + search + "%"
                queryString = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
                args = append(args, search, search, limit)
            }
    }
    rows, err := DB.Query(queryString, args...)
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


func GetTask(id string) (*Task, error) {
    var task Task
    err := DB.QueryRow(
        `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`,
        id,
    ).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

    if err != nil {
        return nil, err
    }
    return &task, nil
}


func UpdateTask(task *Task) error {
    query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
    res, err := DB.Exec(query,  
                                sql.Named("date", task.Date),
								sql.Named("title", task.Title),
								sql.Named("comment", task.Comment),
								sql.Named("repeat", task.Repeat),
                                sql.Named("id", task.ID),
    )
    if err != nil {
        return err
    }
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return sql.ErrNoRows
    }
    return nil
} 


func DeleteTask(id string) error {
    result, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
    if err != nil {
        return err
    }

    count, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return sql.ErrNoRows
    }
    return nil
}
