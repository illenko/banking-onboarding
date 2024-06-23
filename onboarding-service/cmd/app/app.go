package app

import (
	"log"

	"github.com/illenko/onboarding-service/internal/configuration"
	"github.com/illenko/onboarding-service/internal/handler"
	"github.com/illenko/onboarding-service/internal/server"
	"github.com/illenko/onboarding-service/internal/service"
	"github.com/illenko/onboarding-service/internal/worker"
	"go.temporal.io/sdk/client"
)

func Run() {
	configuration.LoadEnv()
	go worker.Run()

	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	defer temporalClient.Close()

	onboardingService := service.NewOnboardingService(temporalClient)
	onboardingHandler := handler.NewOnboardingHandler(onboardingService)

	err = server.New(onboardingHandler).Run(":" + configuration.Get("SERVER_PORT"))
	if err != nil {
		log.Fatalln("Unable to start the server:", err)
		return
	}
}
