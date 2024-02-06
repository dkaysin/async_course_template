package main

import (
	global "async_course/main"
	database "async_course/main/internal/database"
	reader "async_course/main/internal/event_reader"
	writer "async_course/main/internal/event_writer"
	http "async_course/main/internal/http_handler"
	service "async_course/main/internal/service"
	"log/slog"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

// env vars
const (
	listenAddressEnvVar  = "LISTEN_ADDRESS"
	defaultListenAddress = ":4080"

	kafkaBrokersEnvVar        = "KAFKA_BROKERS"
	defaultKafkaBrokersEnvVar = "localhost:9092"

	pgConnStringEnvVar        = "PG_CONN_STRING"
	defaultPgConnStringEnvVar = "postgres://dkaysin:dkaysin@127.0.0.1:5432/async_course_test"
)

func main() {

	// set config
	config := viper.New()
	config.SetEnvPrefix("MAIN")
	config.AutomaticEnv()
	config.SetDefault(listenAddressEnvVar, defaultListenAddress)
	config.SetDefault(kafkaBrokersEnvVar, defaultKafkaBrokersEnvVar)
	config.SetDefault(pgConnStringEnvVar, defaultPgConnStringEnvVar)

	// set database
	db, err := database.NewDatabase(config.GetString(pgConnStringEnvVar))
	if err != nil {
		slog.Error("fatal error while initializing database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// set event writer
	brokers := strings.Split(config.GetString(kafkaBrokersEnvVar), ",")
	ew := writer.NewEventWriter(brokers)
	defer ew.Close()

	// set service
	s := service.NewService(config, db, ew)
	s.ScheduleSendMessages() // TODO: testing

	// set event reader
	er := reader.NewEventReader(s)
	er.StartReaders(brokers, global.KafkaConsumerGroupID)

	// set http handler
	h := http.NewHandler(config, s)

	// set server and API
	e := echo.New()
	api := e.Group("/api")
	h.RegisterAPI(api)

	// set echo logger
	e.Logger.SetPrefix("main")
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}",` +
			`"error_code":"${header:x-hoop-error-code}"}` +
			"\n",
	}))

	// start server
	e.Logger.Fatal(e.Start(config.GetString(listenAddressEnvVar)))
}
