package files

import (
	config "file_server/config"
	model "file_server/model"
	view "file_server/view"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// Get user's files
// @Tags files
// @Description Get user's files
// @Accept  json
// @Produce  json
// @Success 200 {object} view.Response{payload=[]view.FileView}
// @Failure 400,401,406,413,500 {object} view.Response
// @Failure default {object} view.Response
// @Router /files/user-files [get]
// @Security JWT
func GetUserFiles(ctx echo.Context) error {
	claims := ctx.Get("user").(*jwt.Token).Claims.(*view.JwtCustomClaims)
	db := config.GetDB()
	files := []model.File{}

	log.Println(claims.User.ID)
	result := db.Find(&files, "user_id = ?", claims.User.ID)

	if result.Error != nil {
		resp := &view.Response{
			Success: true,
			Message: "Internal Server Error",
			Payload: nil,
		}
		// close db
		config.CloseDB(db).Close()

		return view.ApiView(http.StatusInternalServerError, ctx, resp)
	}

	var formatted_files []view.FileView
	for _, file := range files {
		formatted_files = append(formatted_files, view.FileModelToView(file))
	}

	resp := &view.Response{
		Success: true,
		Message: "Success",
		Payload: formatted_files,
	}

	config.CloseDB(db).Close()

	return view.ApiView(http.StatusOK, ctx, resp)
}
