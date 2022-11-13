package app

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func (app *AppImpl) InitializeDB() error {

	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		app.Config.POSTGRES_HOST, app.Config.POSTGRES_PORT,
		app.Config.POSTGRES_USERNAME, app.Config.POSTGRES_PASSWORD,
		app.Config.POSTGRES_DB)

	DB, err := sql.Open("postgres", dns)
	if err != nil {
		app.ErrorLog.Fatalf("error opening postgres: %s", err.Error())
		return err
	}
	// defer DB.Close()

	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("DB CONNECTED")

	app.DB = DB

	return nil

}

func (app *AppImpl) CheckDB() bool {
	var (
		err    error
		exists = true
		id     int8
	)

	err = app.DB.QueryRow("SELECT id FROM profanities WHERE id = 1").Scan(&id)

	if err != nil || id != 1 {
		exists = false
	}

	return exists
}
