package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"tgstat/internal/config"
	"tgstat/internal/handler/auth"
	channelhandler "tgstat/internal/handler/channel"
	posthandler "tgstat/internal/handler/post"
	userhandler "tgstat/internal/handler/user"
	"tgstat/internal/middleware"
	"tgstat/internal/pkg/client"
	"tgstat/internal/pkg/pg"
	"tgstat/internal/repository"
	"tgstat/internal/service"
	"tgstat/internal/usecase/authorization"
	"tgstat/internal/usecase/bot"
	"tgstat/internal/usecase/channel"
	"tgstat/internal/usecase/post"
	"tgstat/internal/usecase/user"
	"tgstat/openapi"
	"time"
)

const jwtExpireTime time.Duration = time.Hour * 6

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	db, err := pg.NewConnection(cfg.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	reqClient := client.NewClient(cfg.BaseUrl, cfg.Token, &http.Client{})
	userRepo := repository.NewUserRepository(db)
	channelRepo := repository.NewChannelRepository(db)
	postRepo := repository.NewPostRepository(db)

	channelService := service.NewChannelService(reqClient)
	postService := service.NewPostService(reqClient)

	authUseCase := authorization.NewUseCase(userRepo, cfg.JWTSigningKey, jwtExpireTime)
	userUseCase := user.NewUseCase(userRepo, authUseCase)
	postUseCase := post.NewUseCase(channelRepo, postRepo, postService)
	go postUseCase.RunSyncWorker(context.Background())
	channelUseCase := channel.NewUseCase(channelRepo, channelService, postUseCase)

	authHandler := auth.NewHandler(authUseCase)
	userHandler := userhandler.NewHandler(userUseCase)
	channelHandler := channelhandler.NewHandler(channelUseCase, postUseCase)
	postHandler := posthandler.NewHandler(postUseCase)

	tgBot, err := tgbotapi.NewBotAPI("6006223443:AAF8l2ICQmuE9BiZ2Vfzgsh9U7O6HO2rImo")
	if err != nil {
		log.Panic(err)
	}
	botUC := bot.New(tgBot, channelRepo, postRepo)
	go botUC.Run()

	mid := middleware.NewMiddleware(authUseCase)

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	s := http.FS(openapi.FS)
	r.StaticFS("/docs", s)

	apiRouter := r.Group("/api")

	authRouter := apiRouter.Group("/auth")
	authRouter.POST("/sign-in", authHandler.SignIn)

	userRouter := apiRouter.Group("/user")
	userRouter.POST("/create", userHandler.Create)

	channelRouter := apiRouter.Group("/channel").Use(mid.JWT())
	//channelRouter := apiRouter.Group("/channel")
	channelRouter.POST("", channelHandler.CreateChannel)
	channelRouter.GET("", channelHandler.GetAllChannels)
	channelRouter.GET("/:id", channelHandler.GetChannelByID)

	postRouter := apiRouter.Group("/post")
	postRouter.GET("", postHandler.GetAllPost)
	postRouter.GET("/week", postHandler.GetAllForWeek)
	postRouter.GET("/frequency/:channelID", postHandler.GetFrequencyWords)

	if err := r.Run(cfg.Address); err != nil {
		log.Fatal(err)
	}
}
