package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type Task struct {
	ID        int
	Task      string
	Status    int
	Created   time.Time
	Completed *time.Time // Pointer to handle NULL values
}

func main() {
	db, err := sql.Open("mysql", "root:ayush@tcp(127.0.0.1:3306)/todolist?parseTime=true")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	task := flag.String("add", "", "adds task")
	status := flag.Bool("done", false, "mark task as completed")
	id := flag.Int("id", 0, "id of task")
	list := flag.Bool("list", false, "list all tasks")
	remove := flag.Bool("delete", false, "delete task")
	flag.Parse()

	// Adding a task
	if *task != "" {
		stmt, err := db.Prepare("INSERT INTO tasks (task, status, created) VALUES (?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(*task, 0, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Task added successfully!")
	}

	// Updating task status
	if *status {
		if *id == 0 {
			log.Fatal("Please provide a valid task ID using -id flag")
		}

		stmt, err := db.Prepare("UPDATE tasks SET status=?, completed=? WHERE ID=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(1, time.Now(), *id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Task updated successfully!")
	}

	// Listing tasks
	if *list {
		rows, err := db.Query("SELECT ID, task, status, created, completed FROM tasks")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var t Task
			var created time.Time
			var completed sql.NullTime

			err := rows.Scan(
				&t.ID,
				&t.Task,
				&t.Status,
				&created,
				&completed,
			)
			if err != nil {
				log.Fatal(err)
			}

			t.Created = created
			if completed.Valid {
				t.Completed = &completed.Time
			}

			tasks = append(tasks, t)
		}

		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
		} else {
			fmt.Println("Tasks List:")
			for _, t := range tasks {
				completedTime := "Not completed"
				if t.Completed != nil {
					completedTime = t.Completed.Format(time.RFC822)
				}
				fmt.Printf(
					"ID: %d, Task: %s, Status: %d, Created: %s, Completed: %s\n",
					t.ID,
					t.Task,
					t.Status,
					t.Created.Format(time.RFC822),
					completedTime,
				)
			}
		}
	}

	if *remove {
		if *id == 0 {
			log.Fatal("Please provide a valid task ID using -id flag")
		}

		stmt, err := db.Prepare("DELETE FROM tasks WHERE ID=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(*id)

		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Task removed successfully!")
		}
	}
}
