package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rfashwal/scs-dispatcher/internal"
	"github.com/rfashwal/scs-dispatcher/internal/dto"
)

type SensorReadingRequest struct {
	SensorID string
	Type     string
	Value    interface{}
}

func NewRouter(s internal.Service) *gin.Engine {
	r := gin.Default()

	r.POST("/dispatch/sensordata", publishSensorDataHandler(s))

	return r
}

func publishSensorDataHandler(s internal.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		var sensorReadingRequest dto.SensorReadingRequest

		if err := c.Bind(&sensorReadingRequest); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		err := s.PublishSensorData(sensorReadingRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(err)})
			return
		}

		c.JSON(http.StatusOK, "sensor data published")
	}
}
