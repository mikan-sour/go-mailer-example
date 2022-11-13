package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jedzeins/go-mailer/src/config"
	"github.com/jedzeins/go-mailer/src/service"

	"github.com/gin-gonic/gin"
)

type App interface {
	InitService()
	SendEmail(msg service.Message)

	// handlers
	SendMailHandler(c *gin.Context)
	HealthcheckHandler(c *gin.Context)

	// listen
	ListenForMail()
	RunKafkaListener(config *config.Config)

	StartApp(config *config.Config)
}

type AppImpl struct {
	Config *config.Config
	DB     *sql.DB

	Wait     *sync.WaitGroup
	InfoLog  *log.Logger
	ErrorLog *log.Logger

	KafkaListenerService *service.KafkaListenerImpl
	MailService          *service.MailServiceImpl
	ProfanityService     *service.ProfanityDetectionServiceImpl
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

	profService, err := service.NewProfanityDetectionService(app.DB)
	if err != nil {
		app.ErrorLog.Fatalf("error creating a new ProfanityDetectionService: %s", err.Error())
	}

	app.KafkaListenerService = service.NewKafkaListenerService([]string{
		fmt.Sprintf("%s:%s", app.Config.KAFKA_HOST, app.Config.KAFKA_PORT)})
	app.MailService = service.NewMailService(app.Config, app.Wait)
	app.ProfanityService = profService
	app.InfoLog.Println("app services initialized")
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

	// kafka setup
	go app.RunKafkaListener(app.Config)

	// server setup
	app.RunServer()

	app.InfoLog.Println("Server exiting")
}
