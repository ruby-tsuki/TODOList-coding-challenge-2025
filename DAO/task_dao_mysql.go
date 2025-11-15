package DAO

import (
	"J/model"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type MySQLTaskDAO struct {
	db *sql.DB
}

func NewMySQLTaskDAO(dataSourceName string) (*MySQLTaskDAO, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to mysql:%v", err)
	}

	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}
	return &MySQLTaskDAO{
		db: db,
	}, nil
}

func createTable(db *sql.DB) error {
	SQL := `CREATE TABLE IF NOT EXISTS tasks(
    	id INT AUTO_INCREMENT PRIMARY KEY,
    	title VARCHAR(255) NOT NULL,
    	done BOOLEAN NOT NULL DEFAULT FALSE,
    	create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	dead_line TIMESTAMP,
    	INDEX idx_done (done),
		INDEX idx_create_at (create_at)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err := db.Exec(SQL)
	return err

}

func (dao *MySQLTaskDAO) Create(task *model.Task) error {
	query := `INSERT INTO tasks (title,done,create_at,update_at,dead_line) VALUE (?,?,?,?,?)`
	now := time.Now()
	var deadLine interface{}
	if task.DeadLine.IsZero() {
		deadLine = nil
	} else {
		deadLine = task.DeadLine
	}
	res, err := dao.db.Exec(query, task.Title, task.Done, now, now, deadLine)
	if err != nil {
		return fmt.Errorf("failed to insert task:%v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID:%v", err)
	}
	task.ID = int(id)
	task.CreateAt = now
	task.UpdateAt = now
	return nil
}

/*
根据filter返回查询结果
*/
func (dao *MySQLTaskDAO) GetList(filter model.TaskFilter) ([]*model.Task, error) {
	query := "SELECT id, title, done, create_at, update_at, dead_line FROM tasks WHERE done = ?"
	args := []interface{}{filter.Done}

	//默认按照创建时间排序
	if filter.OrderByDeadline {
		query += " ORDER BY dead_line ASC"
	} else {
		query += " ORDER BY create_at DESC"
	}

	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
	}

	rows, err := dao.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %v", err)
	}
	defer rows.Close()

	tasks := make([]*model.Task, 0)
	for rows.Next() {
		task := &model.Task{}
		var deadLine sql.NullTime
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Done,
			&task.CreateAt,
			&task.UpdateAt,
			&deadLine,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %v", err)
		}

		if deadLine.Valid {
			task.DeadLine = deadLine.Time
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return tasks, nil
}

func (dao *MySQLTaskDAO) Update(ID int, title string, done bool, ddl time.Time) error {
	query := `UPDATE tasks SET update_at = ?,title = ?,done = ?,dead_line = ? WHERE id = ?`
	//检测是否需要更新ddl
	var deadLine interface{}
	if ddl.IsZero() {
		deadLine = nil
	} else {
		deadLine = ddl
	}
	temp := []interface{}{time.Now(), title, done, deadLine, ID}
	res, err := dao.db.Exec(query, temp...)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func (dao *MySQLTaskDAO) Delete(taskID int) error {
	query := "DELETE FROM tasks WHERE id = ?"
	res, err := dao.db.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func (dao *MySQLTaskDAO) Count(done bool) (int, error) {
	query := "SELECT COUNT(*) FROM tasks WHERE done = ?"
	var count int
	err := dao.db.QueryRow(query, done).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks: %v", err)
	}
	return count, nil
}

func (dao *MySQLTaskDAO) Close() error {
	return dao.db.Close()
}
