package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

func (app *AppImpl) RunServer() {
	router := gin.Default()
	router.GET("/healthcheck", app.HealthcheckHandler)
	router.POST("/send", app.SendMailHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.Config.API_PORT),
		Handler: router,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		app.InfoLog.Println("receive interrupt signal")
		if err := server.Close(); err != nil {
			app.ErrorLog.Fatal("Server Close:", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			app.ErrorLog.Println("Server closed under request")
		} else {
			app.ErrorLog.Fatal("Server closed unexpect")
		}
	}
}
