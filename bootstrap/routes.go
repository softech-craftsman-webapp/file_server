package bootstrap

import (
	config "file_server/config"
	controller "file_server/controller"
	files "file_server/controller/files"
	_ "file_server/docs"
	helper "file_server/helper"
	"os"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

/*
	|--------------------------------------------------------------------------
	| Routes and its middleware
	|--------------------------------------------------------------------------
*/
func InitRoutes(app *echo.Echo) {
	helper.MakeDirectoryIfNotExists(os.Getenv("BASE_PATH") + "/" + controller.STATIC_PATH)

	// Access, Refresh Application Routes
	access_route := config.Guard(app)

	// enable validation
	app.Validator = &config.CustomValidator{Validator: validator.New()}

	// Swagger
	app.GET("/openapi/*", echoSwagger.WrapHandler)
	app.GET("/openapi", controller.SwaggerRedirect)

	// Files Controller
	access_route.POST("/files/upload", files.Create)
	app.GET("/files/:id/:filename", files.ServeFile)
	access_route.GET("/files/user-files", files.GetUserFiles)
}
