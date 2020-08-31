package main

import (
	"log"
	"time"

	cors "github.com/itsjamie/gin-cors"
	"github.com/rfashwal/scs-dispatcher/internal"
	"github.com/rfashwal/scs-dispatcher/internal/config"
	httptransport "github.com/rfashwal/scs-dispatcher/internal/transport/http"
	"github.com/rfashwal/scs-utilities/rabbit"
)

func main() {

	// For Local Testing
	// _ = os.Setenv("SERVICE_NAME", "sensortypes")
	// _ = os.Setenv("HTTP_PORT", "8102")
	// _ = os.Setenv("EUREKA_SERVICE", "http://127.0.0.1:8761")
	// _ = os.Setenv("ROUTING_KEY", "readings")
	// _ = os.Setenv("TEMPRATURE_TOPIC", "temprature")

	conf := config.Config().Manager
	mq, err := rabbit.NewRabbitMQManager(conf.RabbitURL())
	if err != nil {
		log.Fatalf("MQ server init: %s", err)
	}

	pub, err := mq.InitPublisher()
	if err != nil {
		log.Fatalf("MQ.Publisher init: %s", err)
	}
	err = pub.RabbitConnector.DeclareTopicExchange(conf.TemperatureTopic())
	if err != nil {
		log.Fatal("DeclareTopicExchangeerr", err)
	}

	svc, err := internal.NewService(pub, conf)
	if err != nil {
		log.Fatal("service init err", err)
	}

	router, err := httptransport.NewServer(svc)
	if err != nil {
		log.Fatalf("http server init: %s", err)
	}

	manager := config.EurekaManagerInit()
	manager.SendRegistrationOrFail()
	manager.ScheduleHeartBeat(conf.ServiceName(), 10)
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE, OPTIONS",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	err = router.Run(conf.Address())

	if err != nil {
		log.Fatalf("router run: %s", err)
	}
}
