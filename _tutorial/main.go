package main

import (
	"log"
	"net/http"

	"github.com/teamlint/container"
)

func main() {
	container.Build(
		container.Provide(NewServer),   // provide http server
		container.Provide(NewServeMux), // provide http serve mux
		// endpoints
		container.Provide(NewOrderController, container.As(new(Controller))), // provide order controller
		container.Provide(NewUserController, container.As(new(Controller))),  // provide user controller
	)
	log.Println(container.Instance())

	var server *http.Server
	err := container.Extract(&server)
	if err != nil {
		panic(err)
	}

	server.ListenAndServe()
}

// NewServer creates a http server with provided mux as handler.
func NewServer(mux *http.ServeMux) *http.Server {
	return &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
}

// NewServeMux creates a new http serve mux.
func NewServeMux(controllers []Controller) *http.ServeMux {
	mux := &http.ServeMux{}

	for _, controller := range controllers {
		controller.RegisterRoutes(mux)
	}

	return mux
}

// Controller is an interface that can register its routes.
type Controller interface {
	RegisterRoutes(mux *http.ServeMux)
}

// OrderController is a http controller for orders.
type OrderController struct{}

// NewOrderController creates a auth http controller.
func NewOrderController() *OrderController {
	return &OrderController{}
}

// RegisterRoutes is a Controller interface implementation.
func (a *OrderController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/orders", a.RetrieveOrders)
}

// Retrieve loads orders and writes it to the writer.
func (a *OrderController) RetrieveOrders(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("Orders"))
}

// UserController is a http endpoint for a user.
type UserController struct{}

// NewUserController creates a user http endpoint.
func NewUserController() *UserController {
	return &UserController{}
}

// RegisterRoutes is a Controller interface implementation.
func (e *UserController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", e.RetrieveUsers)
}

// Retrieve loads users and writes it using the writer.
func (e *UserController) RetrieveUsers(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("Users"))
}
