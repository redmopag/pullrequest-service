package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/avito/internship/pr-service/internal/config"
	prHandler "github.com/avito/internship/pr-service/internal/handler/pullrequest"
	teamHandler "github.com/avito/internship/pr-service/internal/handler/team"
	userHandler "github.com/avito/internship/pr-service/internal/handler/user"
	prRepo "github.com/avito/internship/pr-service/internal/repository/pull-request"
	teamRepo "github.com/avito/internship/pr-service/internal/repository/team"
	userRepo "github.com/avito/internship/pr-service/internal/repository/user"
	"github.com/avito/internship/pr-service/internal/router"
	"github.com/avito/internship/pr-service/internal/server"
	prService "github.com/avito/internship/pr-service/internal/service/pullrequest"
	teamService "github.com/avito/internship/pr-service/internal/service/team"
	userService "github.com/avito/internship/pr-service/internal/service/user"
	"github.com/avito/internship/pr-service/internal/storage"
)

func main() {
	cfg := config.NewConfig()

	dbPool := storage.InitDBPool(cfg)
	defer dbPool.Close()

	userRepository := userRepo.NewUserRepository(dbPool)
	teamRepository := teamRepo.NewTeamRepository(dbPool)
	pullRequestRepository := prRepo.NewPullRequestRepository(dbPool)

	userService := userService.NewUserService(userRepository)
	teamService := teamService.NewTeamService(teamRepository)
	pullRequestService := prService.NewPullRequestService(pullRequestRepository, userService, teamService)

	userController := userHandler.NewUserController(userService)
	teamController := teamHandler.NewTeamController(teamService)
	pullRequestController := prHandler.NewPullRequestController(pullRequestService)

	mux := router.InitRouter()
	router.SetupUserRoutes(mux, userController)
	router.SetupTeamRoutes(mux, teamController)
	router.SetupPullRequestRoutes(mux, pullRequestController)

	srv := server.NewServer(cfg.Server.Port, mux)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go srv.Start()
	log.Println("service started")

	<-ctx.Done()

	log.Println("shutting down server...")
	srv.GracefulShutdown()
}
