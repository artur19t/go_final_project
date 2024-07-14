package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	Db *sql.DB
}

var DB Repository

func ConnectDB() error {
	db, err := sql.Open("sqlite", "base/scheduler.db")
	if err != nil {
		return err
	}
	DB = Repository{
		Db: db,
	}
	return nil
}

func CheckBase() {
	//appPath := "base"

	dbFile := filepath.Join("base/scheduler.db")
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		fmt.Println(err)
		install = true
	}

	if install {
		fmt.Println("Создание новой базы")
		_, err = os.Create("base/scheduler.db")
		if err != nil {
			fmt.Println(err)
			return
		}

		db, err := sql.Open("sqlite3", "base/scheduler.db")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = db.Exec(`
			CREATE TABLE scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date CHAR(8) NOT NULL DEFAULT "",
				title VARCHAR(256) NOT NULL DEFAULT "",
				comment TEXT NOT NULL DEFAULT "",
				repeat VARCHAR(128) NOT NULL DEFAULT ""
			);
			CREATE INDEX task_date ON scheduler (date);
		`)
		if err != nil {
			fmt.Println(err)
		}
		db.Close()
	}

}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", nil
	}
	slice := strings.Fields(repeat)
	dateT, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}
	if slice[0] == "d" {
		if len(slice) > 1 {
			num, erro := strconv.Atoi(slice[1])
			if erro != nil {
				return "-1", erro
			}
			if num > 400 {
				return "-1", nil
			}
			timeToReturn := dateT.AddDate(0, 0, num)
			g1 := timeToReturn.Before(now)
			for g1 {

				timeToReturn = timeToReturn.AddDate(0, 0, num)
				g1 = now.After(timeToReturn)
			}

			strTime := timeToReturn.Format("20060102")
			return strTime, nil
		}
	}
	if slice[0] == "y" {
		timeToReturn := dateT.AddDate(1, 0, 0)
		g1 := timeToReturn.Before(now)
		for g1 {

			timeToReturn = timeToReturn.AddDate(1, 0, 0)
			g1 = now.After(timeToReturn)
		}

		strTime := timeToReturn.Format("20060102")
		return strTime, nil
	}
	return "-1", nil
}
