package app

import (
	"log/slog"

	"github.com/illenko/onboarding-service/internal/configuration"
	"github.com/illenko/onboarding-service/internal/handler"
	"github.com/illenko/onboarding-service/internal/server"
	"github.com/illenko/onboarding-service/internal/service"
	"github.com/illenko/onboarding-service/internal/worker"
	"go.temporal.io/sdk/client"
)

func Run() {
	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		slog.Error("Unable to create Temporal client:", slog.String("error", err.Error()))
		return
	}
	defer temporalClient.Close()
	configuration.LoadEnv()
	go worker.Run(temporalClient)

	onboardingService := service.NewOnboardingService(temporalClient)
	onboardingHandler := handler.NewOnboardingHandler(onboardingService)

	err = server.New(onboardingHandler).Run(":" + configuration.Get("SERVER_PORT"))
	if err != nil {
		slog.Error("Unable to start the server:", slog.String("error", err.Error()))
		return
	}
}
