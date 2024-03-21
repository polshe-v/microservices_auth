package main

import (
	"context"

	"go.uber.org/zap"

	"github.com/polshe-v/microservices_auth/internal/app"
	"github.com/polshe-v/microservices_common/pkg/logger"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Fatal("failed to init app: ", zap.Error(err))
	}

	err = a.Run()
	if err != nil {
		logger.Fatal("failed to run app: ", zap.Error(err))
	}
}
