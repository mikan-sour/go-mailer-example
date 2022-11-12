package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (app *AppImpl) HealthcheckHandler(c *gin.Context) {
	time := time.Now().String()
	c.JSON(http.StatusOK, gin.H{"status": "ok", "timeNow": time})
}
