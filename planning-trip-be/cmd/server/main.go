package main

import (
	"log"

	"planning-trip-be/internal/config"
	"planning-trip-be/internal/database"
	"planning-trip-be/internal/middleware"
	arepo "planning-trip-be/internal/repository/album"
	crepo "planning-trip-be/internal/repository/comment"
	nrepo "planning-trip-be/internal/repository/notification"
	phrepo "planning-trip-be/internal/repository/photo"
	ptrepo "planning-trip-be/internal/repository/phototag"
	plrepo "planning-trip-be/internal/repository/place"
	rrepo "planning-trip-be/internal/repository/reaction"
	trepo "planning-trip-be/internal/repository/trip"
	tmrepo "planning-trip-be/internal/repository/tripmember"
	tprepo "planning-trip-be/internal/repository/tripplace"
	tsrepo "planning-trip-be/internal/repository/tripscheduleitem"
	urepo "planning-trip-be/internal/repository/user"
	asvc "planning-trip-be/internal/service/album"
	authsvc "planning-trip-be/internal/service/auth"
	csvc "planning-trip-be/internal/service/comment"
	hsvc "planning-trip-be/internal/service/health"
	nsvc "planning-trip-be/internal/service/notification"
	phsvc "planning-trip-be/internal/service/photo"
	ptsvc "planning-trip-be/internal/service/phototag"
	plsvc "planning-trip-be/internal/service/place"
	rsvc "planning-trip-be/internal/service/reaction"
	tsvc "planning-trip-be/internal/service/trip"
	tmsvc "planning-trip-be/internal/service/tripmember"
	tpsvc "planning-trip-be/internal/service/tripplace"
	tssvc "planning-trip-be/internal/service/tripscheduleitem"
	uploadsvc "planning-trip-be/internal/service/upload"
	usvc "planning-trip-be/internal/service/user"
	httpalbum "planning-trip-be/internal/transport/http/album"
	httpauth "planning-trip-be/internal/transport/http/auth"
	httpcomment "planning-trip-be/internal/transport/http/comment"
	httphealth "planning-trip-be/internal/transport/http/health"
	httpnotification "planning-trip-be/internal/transport/http/notification"
	httpphoto "planning-trip-be/internal/transport/http/photo"
	httpphototag "planning-trip-be/internal/transport/http/phototag"
	httpplace "planning-trip-be/internal/transport/http/place"
	httpreaction "planning-trip-be/internal/transport/http/reaction"
	httptrip "planning-trip-be/internal/transport/http/trip"
	httptripmember "planning-trip-be/internal/transport/http/tripmember"
	httptripplace "planning-trip-be/internal/transport/http/tripplace"
	httptripscheduleitem "planning-trip-be/internal/transport/http/tripscheduleitem"
	httpupload "planning-trip-be/internal/transport/http/upload"
	httpuser "planning-trip-be/internal/transport/http/user"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	env, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	db, err := database.OpenPostgres(env.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("sql database not available: %v", err)
	}
	defer func() {
		if closeErr := sqlDB.Close(); closeErr != nil {
			log.Printf("database close error: %v", closeErr)
		}
	}()

	svc := hsvc.NewService()

	// Initialize repositories
	albumRepo := arepo.NewRepository(db)
	commentRepo := crepo.NewRepository(db)
	notificationRepo := nrepo.NewRepository(db)
	photoRepo := phrepo.NewRepository(db)
	photoTagRepo := ptrepo.NewRepository(db)
	placeRepo := plrepo.NewRepository(db)
	reactionRepo := rrepo.NewRepository(db)
	tripRepo := trepo.NewRepository(db)
	tripMemberRepo := tmrepo.NewRepository(db)
	tripPlaceRepo := tprepo.NewRepository(db)
	tripScheduleRepo := tsrepo.NewRepository(db)
	userRepo := urepo.NewRepository(db)

	// Initialize services
	albumService := asvc.NewService(albumRepo)
	commentService := csvc.NewService(commentRepo)
	notificationService := nsvc.NewService(notificationRepo)
	photoService := phsvc.NewService(photoRepo)
	photoTagService := ptsvc.NewService(photoTagRepo)
	placeService := plsvc.NewService(placeRepo)
	reactionService := rsvc.NewService(reactionRepo)
	tripService := tsvc.NewService(tripRepo)
	tripMemberService := tmsvc.NewService(tripMemberRepo)
	tripPlaceService := tpsvc.NewService(tripPlaceRepo)
	tripScheduleService := tssvc.NewService(tripScheduleRepo)
	uploadService := uploadsvc.NewService(
		env.CloudinaryCloudName,
		env.CloudinaryAPIKey,
		env.CloudinaryAPISecret,
		env.CloudinaryUploadPath,
	)
	userService := usvc.NewService(userRepo)
	authService := authsvc.NewService(userRepo, env.AuthSecret)

	router := gin.New()
	router.Use(gin.Recovery(), middleware.Logging(), middleware.CORS())
	authMiddleware := middleware.AuthRequired(env.AuthSecret)

	// Health check endpoint
	router.GET("/health", httphealth.Handler(svc))

	// Register routes
	httpalbum.RegisterRoutes(router, albumService)
	httpauth.RegisterRoutes(router, authService)
	httpcomment.RegisterRoutes(router, commentService)
	httpnotification.RegisterRoutes(router, notificationService)
	httpphoto.RegisterRoutes(router, photoService)
	httpphototag.RegisterRoutes(router, photoTagService)
	httpplace.RegisterRoutes(router, placeService)
	httpreaction.RegisterRoutes(router, reactionService)
	httptrip.RegisterRoutes(router, tripService, authMiddleware)
	httptripmember.RegisterRoutes(router, tripMemberService)
	httptripplace.RegisterRoutes(router, tripPlaceService)
	httptripscheduleitem.RegisterRoutes(router, tripScheduleService)
	httpupload.RegisterRoutes(router, uploadService, authMiddleware)
	httpuser.RegisterRoutes(router, userService)

	log.Printf("Starting planning-trip-be on port %s", env.Port)
	if err := router.Run(":" + env.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
