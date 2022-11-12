package app

import (
	"net/http"

	"github.com/jedzeins/go-mailer/src/service"

	"github.com/gin-gonic/gin"
)

func (app *AppImpl) SendMailHandler(c *gin.Context) {
	var msg service.Message

	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "nok", "error": err})
		return
	}

	app.SendEmail(msg)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "result": msg})

}
