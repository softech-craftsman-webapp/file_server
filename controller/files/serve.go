package files

import (
	config "file_server/config"
	controller "file_server/controller"
	helper "file_server/helper"
	model "file_server/model"
	view "file_server/view"
	"fmt"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
)

/*
   |--------------------------------------------------------------------------
   | Serve files
   | @Param filename
   | @Query download => @url?download => download file
   |--------------------------------------------------------------------------
*/
func ServeFile(ctx echo.Context) error {
	db := config.GetDB()
	f := &model.File{}

	filename := ctx.Param("filename")
	id := ctx.Param("id")
	fileDir := path.Clean(fmt.Sprintf("%v/%v", controller.STATIC_PATH, filename))

	/*
		|--------------------------------------------------------------------------
		| File Database
		|--------------------------------------------------------------------------
	*/
	file_result := db.First(&f, "name = ? AND id = ?", filename, id)
	if file_result.Error != nil {
		resp := &view.Response{
			Success: false,
			Message: "Not found",
			Payload: nil,
		}

		// close db
		config.CloseDB(db).Close()

		return view.ApiView(http.StatusNotFound, ctx, resp)
	}

	// File types
	if !helper.Contains(controller.MIME_TYPES, f.Type) {
		resp := &view.Response{
			Success: false,
			Message: "File is not permitted to serve",
			Payload: nil,
		}

		return view.ApiView(http.StatusForbidden, ctx, resp)
	}

	config.CloseDB(db).Close()

	// Download
	if download := ctx.QueryParam("download"); download == "" {
		ctx.Attachment(fileDir, f.Name)
	}

	// Serve file
	return ctx.Inline(fileDir, f.Name)
}
