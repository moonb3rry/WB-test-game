package app

import (
	"WB_game/config"
	"WB_game/internal/repository"
	"WB_game/internal/service"
	"WB_game/internal/transport/http"
	"WB_game/pkg/httpserver"
	"WB_game/pkg/postgres"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	pg, err := postgres.New(postgres.GetConnString(&cfg.Db), postgres.MaxPoolSize(cfg.Db.MaxPoolSize))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pg.Close()

	err = pg.Pool.Ping(context.Background())
	if err != nil {
		fmt.Print(err)
	}

	userRepo := repository.NewUserRepository(pg)
	taskRepo := repository.NewTaskRepository(pg)
	gameRepo := repository.NewGameRepository(pg)
	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)
	gameService := service.NewGameService(gameRepo, userRepo, taskRepo)
	httpController := http.NewController(userService, taskService, gameService)

	if err != nil {
		log.Fatal("Controller initialization problem: %v", err)
	}

	httpServer := httpserver.New(httpController,
		httpserver.Port(cfg.HttpServer.Addr),
		httpserver.ReadTimeout(cfg.HttpServer.ReadTimeout),
		httpserver.WriteTimeout(cfg.HttpServer.WriteTimeout),
		httpserver.ShutdownTimeout(cfg.HttpServer.ShutdownTimeout),
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	log.Printf("Running app:%v version:%v", cfg.App.Name, cfg.App.Version)

	select {
	case s := <-interrupt:
		log.Printf("App running signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Fatal(fmt.Errorf("App HTTP server notify: %v", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		log.Fatal(fmt.Errorf("App HTTP server shutdown: %v", err))
	}
	fmt.Printf("server down")
}
