package main

import (
	"github.com/jedzeins/go-mailer/src/app"
	"github.com/jedzeins/go-mailer/src/config"
)

func main() {
	config, err := config.New()
	if err != nil {
		panic(err)
	}

	application := app.New(config)

	application.InitializeDB()

	application.InitService()

	go application.ListenForMail()

	application.StartApp(application.Config)

}
