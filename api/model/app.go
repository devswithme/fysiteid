package model

import (
	"fybe/controller"
	"fybe/repository"
	"fybe/service"
)

type App struct {
	Repository struct {
		User       repository.UserRepository
		Ticket     repository.TicketRepository
		Registrant repository.RegistrantRepository
	}

	Service struct {
		User       service.UserService
		Auth       service.AuthService
		Ticket     service.TicketService
		Registrant service.RegistrantService
	}

	Controller struct {
		User       controller.UserController
		Auth       controller.AuthController
		Ticket     controller.TicketController
		Registrant controller.RegistrantController
	}
}
