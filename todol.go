package main

import (
	"fmt"
	"github.com/mxk/go-sqlite/sqlite3"
	"log"
	"os"
	"os/user"
	"time"
)

func main() {
	// config
	prog_version := "0.1"
	currentTime := time.Now().Format("02-01-2006, 15:04:05")
	args := os.Args[1:]
	// -----

	// sqlite
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	c, _ := sqlite3.Open(usr.HomeDir + "/.config/todol.db")
	// create a base table
	query := "CREATE TABLE notes(id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, date VARCHAR(21), text TEXT)"
	c.Exec(query)
	// -----

	if len(args) == 0 {
		fmt.Println("Todol v"+prog_version+" by CubexX\n\n",
			"Usage:\n",
			// help
			"\033[36m add\033[32m <text>\033[0m - add a note to your Todol\n",
			"\033[36m show\033[32m <all/id>\033[0m - show your notes/note by id\n",
			//"\033[36m edit\033[32m <id> <new text>\033[0m - edit your note by id\n",
			"\033[36m del\033[32m <all/id>\033[0m - delete all notes/by id")
	} else {
		switch args[0] {
		// add note
		case "add":
			if len(args[1:]) == 0 { // if not exists text of note
				fmt.Println("Use\033[36m todol add \033[32m<text>\033[0m")
			} else {
				fmt.Println("\033[32m\"" + args[1] + "\"\033[0m added to your Todol")
				// sqlite add note
				values := sqlite3.NamedArgs{"$text": args[1], "$date": currentTime}
				c.Exec("INSERT INTO notes(text, date) VALUES($text, $date)", values)

			}
		// show notes
		case "show":
			if len(args[1:]) == 0 { // if not exists id
				fmt.Println("Use\033[36m todol show \033[32m<all/id>\033[0m")
			} else {
				// show all notes
				if args[1] == "all" {
					sql := "SELECT * FROM notes"
					i := make(sqlite3.RowMap)
					d := make(sqlite3.RowMap)
					t := make(sqlite3.RowMap)
					for s, err := c.Query(sql); err == nil; err = s.Next() {
						var (
							id   int
							date string
							text string
						)
						s.Scan(&id, i)
						s.Scan(&date, d)
						s.Scan(&text, t)
						fmt.Println("\033[31mID:", id, "\033[34m", d["date"], "\033[0m\n", t["text"])
					}
				} else { // show note by id
					sql := "SELECT * FROM notes WHERE id=" + args[1]
					i := make(sqlite3.RowMap)
					d := make(sqlite3.RowMap)
					t := make(sqlite3.RowMap)
					for s, err := c.Query(sql); err == nil; err = s.Next() {
						var (
							id   int
							date string
							text string
						)
						s.Scan(&id, i)
						s.Scan(&date, d)
						s.Scan(&text, t)
						if len(string(id)) != 0 {
							fmt.Println("\033[31mID:", id, "\033[34m", d["date"], "\033[0m\n", t["text"])
						} else {
							fmt.Println("\033[31mNote with this ID doesn't exist!\033[0m")
						}
					}
				}
			}
		case "edit":
			if len(args[1:]) == 0 { // if not exists id
				fmt.Println("Use\033[36m todol edit\033[32m <id> \033[0m<text>")
			} else {
				if len(args[2:]) == 0 { // if not exists text
					fmt.Println("Use\033[36m todol edit\033[0m <id> \033[32m<text>\033[0m")
				} else {
					q := "UPDATE notes SET text='" + args[2] + "' WHERE id=" + args[1]
					c.Exec(q)
					fmt.Println("\033[32mNote with ID \033[31m" + args[1] + "\033[32m was updated\033[0m")
				}
			}
		case "del":
			if len(args[1:]) == 0 { // if not exists id
				fmt.Println("Use\033[36m todol show \033[32m<all/id>\033[0m")
			} else {
				// delete all notes
				if args[1] == "all" {
					q := "DELETE FROM notes"
					c.Exec(q)
					fmt.Println("\033[31mAll notes was deleted\033[0m")
				} else { // show note by id
					q := "DELETE FROM notes WHERE id=" + args[1]
					c.Exec(q)
					fmt.Println("\033[32mNote with ID \033[31m" + args[1] + "\033[32m was deleted\033[0m")
				}
			}
		// if not exists command
		default:
			fmt.Println("\033[31mInvalid command!\033[0m")
		}
	}
}
