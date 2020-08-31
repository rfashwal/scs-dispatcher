package internal

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rfashwal/scs-dispatcher/internal/dto"
	"github.com/rfashwal/scs-utilities/config"
	"github.com/rfashwal/scs-utilities/rabbit/domain"
	"github.com/rfashwal/scs-utilities/rabbit/publishing"
)

type Service interface {
	PublishSensorData(req dto.SensorReadingRequest) error
}

func NewService(publisher *publishing.Publisher, conf config.Manager) (Service, error) {
	return service{publisher: publisher}, nil
}

type service struct {
	publisher *publishing.Publisher
	conf      config.Manager
}

func (s service) PublishSensorData(req dto.SensorReadingRequest) error {

	switch req.SensorType {
	case "temprature":
		parsedValue, err := strconv.ParseFloat(req.Value, 64)
		if err != nil {
			return err
		}
		msg := domain.TemperatureMeasurement{
			ProcessId:   uuid.New().String(),
			PublishedOn: time.Now(),
			SensorId:    req.SensorId,
			Value:       parsedValue,
			Service:     s.conf.ServiceName(),
			RoomId:      req.RoomId,
		}
		encodedMsg, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		return s.publisher.Publish(s.conf.TemperatureTopic(), s.conf.ReadingsRoutingKey(), string(encodedMsg))
	}
	return nil
}
