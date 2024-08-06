package main

import (
	"context"
	"fmt"
	"net/http"
	_ "tgs-automation/features/domain-api/docs"
	jwttoken "tgs-automation/internal/jwt_token"
	middleware "tgs-automation/internal/middleware"
	"tgs-automation/internal/opentelemetry"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/cloudflare"
	"tgs-automation/pkg/namecheap"

	"github.com/gin-gonic/gin"
	"github.com/iris-contrib/swagger/swaggerFiles"

	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	ctx := context.Background()
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	config := util.GetConfig()
	fmt.Println("config.JaegerCollectorUrl:", config.JaegerCollectorUrl)
	opentelemetry.InitTracerProvider(ctx, config.JaegerCollectorUrl, "domain-api", "0.1.0", "prod")
	router.Use(middleware.TraceMiddleware("domain-api"))
	natsSvc, err := util.NewNatsPublisher(config.NatsUrl)
	defer natsSvc.Close()

	if err != nil {
		fmt.Println("cannot get nats svc.")
	}

	namecheapSvc := namecheap.New(config)
	cloudflareSvc := cloudflare.NewClloudflare(config)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1", "10.139.0.0/16"})
	router.GET("/healthz", healthCheckHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/token", jwttoken.TokenHandler)
	router.GET("/nameserver", middleware.AuthMiddleware(), GetNameServerHandler(cloudflareSvc, natsSvc))
	router.PUT("/nameserver", middleware.AuthMiddleware(), UpdateNameServerHandler(namecheapSvc, natsSvc))
	router.GET("/domain/price", middleware.AuthMiddleware(), GetDomainPriceHandler(namecheapSvc, natsSvc))
	router.POST("/domain", middleware.AuthMiddleware(), CreateDomainHandler(namecheapSvc, natsSvc))
	err = server.ListenAndServe()

	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "up"})
}
