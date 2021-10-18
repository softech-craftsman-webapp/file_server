package files

import (
	"bytes"
	config "file_server/config"
	controller "file_server/controller"
	helper "file_server/helper"
	model "file_server/model"
	view "file_server/view"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Create a new file
// @Tags files
// @Description Create a new file
// @Accept  json
// @Produce  json
// @Param   file formData file true  "File"
// @Success 201 {object} view.Response{payload=view.FileView}
// @Failure 400,401,406,413,500 {object} view.Response
// @Failure default {object} view.Response
// @Router /files/upload [post]
// @Security JWT
func Create(ctx echo.Context) error {
	fileID := uuid.New().String()
	claims := ctx.Get("user").(*jwt.Token).Claims.(*view.JwtCustomClaims)

	// Source
	file, err := ctx.FormFile("file")
	file.Filename = fileID + filepath.Ext(file.Filename)

	/*
	   |--------------------------------------------------------------------------
	   | Check if file is valid
	   |--------------------------------------------------------------------------
	*/
	if err != nil {
		resp := &view.Response{
			Success: false,
			Message: err.Error(),
			Payload: nil,
		}

		return view.ApiView(http.StatusBadRequest, ctx, resp)
	}

	/*
	   |--------------------------------------------------------------------------
	   | Check file size
	   |--------------------------------------------------------------------------
	*/
	if file.Size > controller.MAX_UPLOAD_SIZE {
		resp := &view.Response{
			Success: false,
			Message: "File is bigger",
			Payload: nil,
		}

		return view.ApiView(http.StatusRequestEntityTooLarge, ctx, resp)
	}

	/*
	   |--------------------------------------------------------------------------
	   | Check file type
	   |--------------------------------------------------------------------------
	*/
	src, err := file.Open()

	if err != nil {
		resp := &view.Response{
			Success: false,
			Message: err.Error(),
			Payload: nil,
		}

		return view.ApiView(http.StatusInternalServerError, ctx, resp)
	}

	defer src.Close()

	buff := bytes.NewBuffer(nil)
	if _, err := io.Copy(buff, src); err != nil {
		resp := &view.Response{
			Success: false,
			Message: err.Error(),
			Payload: nil,
		}

		return view.ApiView(http.StatusInternalServerError, ctx, resp)
	}

	filetype := http.DetectContentType(buff.Bytes())

	// File types
	if !helper.Contains(controller.MIME_TYPES, filetype) {
		resp := &view.Response{
			Success: false,
			Message: "File is not permitted to upload",
			Payload: nil,
		}

		return view.ApiView(http.StatusNotAcceptable, ctx, resp)
	}

	/*
	   |--------------------------------------------------------------------------
	   | Destionation
	   |--------------------------------------------------------------------------
	*/
	dst, err := os.Create(os.Getenv("BASE_PATH") + "/" + controller.STATIC_PATH + "/" + file.Filename)

	if err != nil {
		resp := &view.Response{
			Success: false,
			Message: err.Error(),
			Payload: nil,
		}

		return view.ApiView(http.StatusInternalServerError, ctx, resp)
	}

	defer dst.Close()

	/*
	   |--------------------------------------------------------------------------
	   | Copy
	   |--------------------------------------------------------------------------
	*/

	src, err = file.Open()

	if err != nil {
		resp := &view.Response{
			Success: false,
			Message: err.Error(),
			Payload: nil,
		}

		return view.ApiView(http.StatusInternalServerError, ctx, resp)
	}

	defer src.Close()

	if _, err = io.Copy(dst, src); err != nil {
		resp := &view.Response{
			Success: false,
			Message: err.Error(),
			Payload: nil,
		}

		return view.ApiView(http.StatusInternalServerError, ctx, resp)
	}

	/*
	   |--------------------------------------------------------------------------
	   | Success
	   |--------------------------------------------------------------------------
	*/
	f := &model.File{
		Name:   file.Filename,
		UserID: claims.User.ID,
		Size:   file.Size,
		Type:   filetype,
	}

	db := config.GetDB()
	result := db.Create(&f)

	resp := &view.Response{
		Success: true,
		Message: "Success",
		Payload: &view.FileView{
			Name: f.Name,
			Size: helper.ByteCountSI(int64(f.Size)),
			Type: f.Type,
			Url:  fmt.Sprintf("%v/%v/%v/%v", os.Getenv("PUBLIC_URL"), controller.STATIC_PATH, f.ID, f.Name),
		},
	}

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

	config.CloseDB(db).Close()

	return view.ApiView(http.StatusCreated, ctx, resp)
}
