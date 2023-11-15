package http

import (
	"WB_game/internal/middleware"
	"github.com/gorilla/mux"
)

func NewController(us UserService, ts TaskService, gs GameService) *mux.Router {
	controller := newController(us, ts, gs)
	router := mux.NewRouter()

	router.HandleFunc("/register", controller.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/login", controller.LogInUserHandler).Methods("POST")
	router.HandleFunc("/get-tasks", controller.GetTasksHandler).Methods("GET")

	protectedRouter := router.PathPrefix("").Subrouter()
	protectedRouter.Use(middleware.JWTAuthentication)
	protectedRouter.HandleFunc("/me", controller.AboutHandler).Methods("GET")
	protectedRouter.HandleFunc("/tasks", controller.TasksHandler).Methods("GET")
	protectedRouter.HandleFunc("/start", controller.StartGameHandler).Methods("POST")

	return router
}
