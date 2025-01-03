package route

import (
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/handler"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RouteConfig struct {
	App            *gin.Engine
	Log            *logrus.Logger
	Viper          *viper.Viper
	UserHandler    handler.IUserHandler
	AuthMiddleware gin.HandlerFunc
}

func (c *RouteConfig) SetupRoutes() {
	c.App.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello world",
		})
	})

	c.SetupAPIRoutes()
}

func (c *RouteConfig) SetupAPIRoutes() {
	apiRoute := c.App.Group("/api")
	{
		apiRoute.POST("/login", c.UserHandler.Login)
		apiRoute.Use(c.AuthMiddleware)
		{
			apiRoute.GET("/me", func(ctx *gin.Context) {
				ctx.JSON(200, gin.H{
					"message": "Welcome brother",
				})
			})
		}
	}
}

func NewRouteConfig(app *gin.Engine, viper *viper.Viper, log *logrus.Logger) *RouteConfig {
	// factory handlers

	// factory middleware
	authMiddleware := middleware.NewAuth(viper)
	return &RouteConfig{
		App:            app,
		Log:            log,
		Viper:          viper,
		UserHandler:    handler.UserHandlerFactory(log, viper),
		AuthMiddleware: authMiddleware,
	}
}
