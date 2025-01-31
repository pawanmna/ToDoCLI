package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"strings"
	"time"
)

type Task struct {
	ID        int
	Task      string
	Status    int
	Created   time.Time
	Completed *time.Time
}

func printLine(width map[string]int) {
	fmt.Print(" ")
	for _, w := range []string{"id", "task", "status", "created", "completed"} {
		fmt.Print(strings.Repeat(" ", width[w]+2) + " ")
	}
	fmt.Println()
}

func printRow(id, task, status, created, completed string, width map[string]int) {
	fmt.Printf(" %-*s  %-*s  %-*s  %-*s  %-*s \n",
		width["id"], id,
		width["task"], task,
		width["status"], status,
		width["created"], created,
		width["completed"], completed)
}

func main() {
	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatal(err)
	}

	// Connect to SQLite database
	db, err := sql.Open("sqlite", "file:todolist.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			ID INTEGER PRIMARY KEY,
			task TEXT NOT NULL,
			status INTEGER NOT NULL DEFAULT 0,
			created DATETIME NOT NULL,
			completed DATETIME
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	task := flag.String("add", "", "adds task")
	status := flag.Bool("done", false, "mark task as completed")
	id := flag.Int("id", 0, "id of task")
	list := flag.Bool("list", false, "list all tasks")
	remove := flag.Bool("delete", false, "delete task")
	All := flag.Bool("all", false, "list all tasks including completed")
	flag.Parse()

	if *task != "" {
		stmt, err := db.Prepare("INSERT INTO tasks (task, status, created) VALUES (?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(*task, 0, time.Now().UTC())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Task added successfully!")
	}

	if *status {
		if *id == 0 {
			log.Fatal("Please provide a valid task ID using -id flag")
		}
		stmt, err := db.Prepare("UPDATE tasks SET status=?, completed=? WHERE ID=?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(1, time.Now().UTC(), *id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Task updated successfully!")
	}

	if *list {
		var query string
		if *All {
			query = "SELECT ID, task, status, created, completed FROM tasks"
		} else {
			query = "SELECT ID, task, status, created, completed FROM tasks WHERE status = 0"
		}
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var tasks []Task
		width := map[string]int{
			"id":        2,
			"task":      4,
			"status":    7,
			"created":   19,
			"completed": 19,
		}

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
			} else {
				t.Completed = nil
			}
			tasks = append(tasks, t)

			if w := len(fmt.Sprintf("%d", t.ID)); w > width["id"] {
				width["id"] = w
			}
			if w := len(t.Task); w > width["task"] {
				width["task"] = w
			}
		}
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
		} else {
			printLine(width)
			printRow("ID", "Task", "Status", "Created", "Completed", width)
			printLine(width)

			for _, t := range tasks {
				createdIST := t.Created.In(ist)
				var completedTime string
				if t.Completed != nil {
					completedIST := t.Completed.In(ist)
					completedTime = completedIST.Format("02 Jan 06 15:04 IST")
				} else {
					completedTime = "Not completed"
				}

				status := "Pending"
				if t.Status == 1 {
					status = "Done"
				}

				printRow(
					fmt.Sprintf("%d", t.ID),
					t.Task,
					status,
					createdIST.Format("02 Jan 06 15:04 IST"),
					completedTime,
					width,
				)
			}
			printLine(width)
		}
	}

	if *remove {
		if *remove && *All {
			_, err := db.Exec("DELETE FROM tasks")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("All tasks deleted successfully!")
		} else {
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
			}
			fmt.Printf("Task %d deleted successfully!\n", *id)
		}
	}
}
