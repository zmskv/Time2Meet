package http

import (
	"time2meet/internal/application/usecase/batch"
	"time2meet/internal/application/usecase/event"
	"time2meet/internal/application/usecase/report"
	"time2meet/internal/application/usecase/ticket"
	"time2meet/internal/application/usecase/user"
	"time2meet/internal/application/usecase/venue"
	"time2meet/internal/infrastructure/persistence/postgres"
	"time2meet/internal/presentation/http/handler"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "time2meet/docs/swagger"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Dependencies struct {
	DB  *sqlx.DB
	Log *zap.Logger
}

func NewRouter(deps Dependencies) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	userRepo := postgres.NewUserRepo(deps.DB)
	userProfileRepo := postgres.NewUserProfileRepo(deps.DB)
	eventRepo := postgres.NewEventRepo(deps.DB)
	venueRepo := postgres.NewVenueRepo(deps.DB)
	roomRepo := postgres.NewRoomRepo(deps.DB)
	ticketRepo := postgres.NewTicketRepo(deps.DB)
	reportRepo := postgres.NewReportRepo(deps.DB)
	txManager := postgres.NewTxManager(deps.DB, deps.Log)
	auditCtx := postgres.NewAuditContextSetter()
	ticketTx := postgres.NewTicketTxQueries()
	batchImp := postgres.NewBatchImporter(deps.Log)

	userUC := user.New(userRepo, userProfileRepo)
	eventUC := event.New(eventRepo)
	venueUC := venue.New(venueRepo, roomRepo)
	reportUC := report.New(reportRepo)
	purchaseUC := ticket.NewPurchase(txManager, auditCtx, ticketTx)
	ticketUC := ticket.NewTicketUC(ticketRepo)
	validateUC := ticket.NewValidate(txManager, auditCtx, ticketTx)
	batchUC := batch.New(txManager, auditCtx, batchImp)

	userH := handler.NewUserHandler(userUC)
	eventH := handler.NewEventHandler(eventUC)
	venueH := handler.NewVenueHandler(venueUC)
	reportH := handler.NewReportHandler(reportUC)
	ticketH := handler.NewTicketHandler(purchaseUC, ticketUC, validateUC)
	batchH := handler.NewBatchHandler(batchUC)

	api := r.Group("/api/v1")
	{
		api.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.URL("/api/v1/swagger/doc.json"),
		))

		api.GET("/users", userH.List)
		api.POST("/users", userH.Create)
		api.GET("/users/:id", userH.Get)
		api.PUT("/users/:id", userH.Update)
		api.DELETE("/users/:id", userH.Delete)

		api.GET("/events", eventH.List)
		api.POST("/events", eventH.Create)
		api.GET("/events/:id", eventH.Get)
		api.PUT("/events/:id", eventH.Update)
		api.DELETE("/events/:id", eventH.Delete)
		api.POST("/events/:id/cancel", eventH.Cancel)

		api.POST("/venues", venueH.CreateVenue)
		api.GET("/venues", venueH.ListVenues)
		api.GET("/venues/:id", venueH.GetVenue)
		api.PUT("/venues/:id", venueH.UpdateVenue)
		api.DELETE("/venues/:id", venueH.DeleteVenue)
		api.POST("/venues/:id/rooms", venueH.CreateRoom)
		api.GET("/venues/:id/rooms", venueH.ListRooms)

		api.POST("/tickets/purchase", ticketH.Purchase)
		api.GET("/tickets/:id", ticketH.Get)
		api.GET("/tickets", ticketH.ListByBuyer)
		api.PATCH("/tickets/:id/status", ticketH.UpdateStatus)
		api.DELETE("/tickets/:id", ticketH.Delete)
		api.POST("/tickets/:id/validate", ticketH.Validate)

		api.GET("/reports/sales", reportH.Sales)
		api.GET("/reports/attendance", reportH.Attendance)
		api.GET("/analytics/popular-events", reportH.Popular)

		api.POST("/batch/import/users", batchH.ImportUsers)
		api.POST("/batch/import/events", batchH.ImportEvents)
		api.POST("/batch/import/tickets", batchH.ImportTickets)
	}

	return r
}
