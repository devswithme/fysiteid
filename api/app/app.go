package app

import (
	"fybe/controller"
	"fybe/helper"
	"fybe/model"
	"fybe/repository"
	"fybe/service"
	"os"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitApp(db *gorm.DB, logger *zap.Logger, redis helper.RedisHelper) *model.App {
	app := &model.App{}
	upload := helper.NewUploadHelper(os.Getenv("PARENT_UPLOAD"), 2*1024*1024, []string{".png", ".jpeg", ".jpg"}, logger)

	app.Repository.User = repository.NewUserRepository(db, logger)
	app.Repository.Ticket = repository.NewTicketRepository(db, logger)
	app.Repository.Registrant = repository.NewRegistrantRepository(db, logger)

	app.Service.User = service.NewUserService(app.Repository.User, upload, logger)
	app.Service.Auth = service.NewAuthService(app.Repository.User, app.Service.User, redis, logger)
	app.Service.Ticket = service.NewTicketService(app.Repository.Ticket, upload, redis, logger)
	app.Service.Registrant = service.NewRegistrantService(app.Repository.Registrant, app.Repository.Ticket, redis, logger)

	app.Controller.User = controller.NewUserController(app.Service.User, logger)
	app.Controller.Auth = controller.NewAuthController(app.Service.Auth, logger)
	app.Controller.Ticket = controller.NewTicketController(app.Service.Ticket, logger)
	app.Controller.Registrant = controller.NewRegistrantService(app.Service.Registrant, logger)

	return app
}
