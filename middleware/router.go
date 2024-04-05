package middleware

import (
	"net/http"
	"web/config"
	"web/handlers"
	"web/model"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, env model.AppConfig) *gin.Engine {
	// Setup Security Headers and check for valid host
	router.Use(func(c *gin.Context) {
		config.MiddlewareHTTPConfig(c, env)
	})

	// Load HTML Templates
	router.LoadHTMLGlob("./public/views/*.html")

	// Serve static files
	router.Use(static.Serve("/", static.LocalFile("./public/", true)))

	// API Routes
	router.GET(env.APIPath+"/ping", handlers.PingAPIHandler)
	router.GET("/status", handlers.StatusHandler(env), func(c *gin.Context) {
		statuses := c.MustGet("statuses").([]map[string]string)
		c.HTML(http.StatusOK, "status.html", gin.H{
			"services": statuses,
		})
	})

	// Error Handlers
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "public/error/404.html", nil)
	})
	router.Use(handlers.ServerErrorHandler())
	router.NoRoute(handlers.NotFoundHandler)

	return router
}