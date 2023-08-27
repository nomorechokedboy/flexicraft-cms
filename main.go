package main

import (
	apperr "api/src/app_error"
	"api/src/config"
	"api/src/controllers"
	"api/src/entities"
	"api/src/hasher"
	"api/src/repositories"
	"api/src/tokenizer"
	"api/src/usecases"
	"api/src/validator"
	"fmt"
	"net/http"

	_ "api/docs"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type OAuthQuery struct {
	Code     string `query:"code"`
	Provider string `query:"provider"`
}

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:5000
// @BasePath /api/v1
func main() {
	config := config.LoadEnv()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&entities.AuthEntity{}, &Collection{}); err != nil {
		panic("failed to migrate")
	}

	p := hasher.Argon2HashParams{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
	hasher := hasher.Argon2Hasher{Argon2HashParams: p}
	gormAuthRepo := repositories.New(db)
	jwt := tokenizer.Jwt{}
	validator := validator.New()
	authUseCase := usecases.New(hasher, gormAuthRepo, jwt, validator)
	authController := controllers.New(authUseCase)

	oauthUsecase := usecases.NewOAuth(jwt)
	oauthController := controllers.NewOAuthController(oauthUsecase)

	app := echo.New()
	app.Use(middleware.CORS())
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: app.Logger.Output(),
	}))
	app.HTTPErrorHandler = func(err error, c echo.Context) {
		if e, ok := err.(*apperr.AppError); ok {
			c.JSON(e.HTTPCode, e)
			return
		}

		app.DefaultHTTPErrorHandler(err, c)
	}
	testie := Testie{db}

	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.POST("/signup", authController.SignUp)
	v1.POST("/signin", authController.SignIn)
	v1.POST("/test", testie.CreateCollection)
	v1.GET("/test", testie.GetCollection)
	v1.POST("/test-record/:collection", testie.CreateRecord)

	app.GET("/docs/*", echoSwagger.WrapHandler)
	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	app.GET("/api/sessions/oauth/github", oauthController.AuthWithOAuth2)
	app.GET("/login/github", func(c echo.Context) error {
		redirectURL := fmt.Sprintf(
			"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
			config.GitHubClientID,
			config.GitHubOAuthRedirectUrl,
		)

		return c.Redirect(http.StatusMovedPermanently, redirectURL)
	})

	app.Start(":5000")
}
