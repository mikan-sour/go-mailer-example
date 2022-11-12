package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jedzeins/go-mailer/src/config"
	"github.com/jedzeins/go-mailer/src/service"
)

type App interface {
}

type AppImpl struct {
	Config *config.Config

	Wait     *sync.WaitGroup
	InfoLog  *log.Logger
	ErrorLog *log.Logger

	MailService *service.MailServiceImpl
}

func New(config *config.Config) *AppImpl {

	wg := sync.WaitGroup{}

	return &AppImpl{
		Config:   config,
		Wait:     &wg,
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (app *AppImpl) InitService() {
	app.MailService = service.NewMailService(app.Config, app.Wait)
	app.InfoLog.Println("app mail service initialized")
}

func (app *AppImpl) SendEmail(msg service.Message) {
	app.Wait.Add(1)
	app.MailService.MailerChan <- msg
}

func (app *AppImpl) ListenForMail() {
	for {
		select {
		case msg := <-app.MailService.MailerChan:
			go app.MailService.SendMail(msg, app.MailService.ErrorChan)
		case err := <-app.MailService.ErrorChan:
			app.ErrorLog.Println(err)
		case <-app.MailService.DoneChan:
			return
		}
	}
}

func (app *AppImpl) StartApp(config *config.Config) {

	router := gin.Default()
	router.GET("/healthcheck", app.HealthcheckHandler)
	router.POST("/send", app.SendMailHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.API_PORT),
		Handler: router,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("receive interrupt signal")
		if err := server.Close(); err != nil {
			log.Fatal("Server Close:", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect")
		}
	}

	log.Println("Server exiting")
}