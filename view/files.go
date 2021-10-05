package view

import (
	controller "file_server/controller"
	helper "file_server/helper"
	model "file_server/model"
	"fmt"
	"os"
)

type FileView struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Type string `json:"type"`
	Size string `json:"size"`
}

func FileModelToView(file model.File) FileView {
	return FileView{
		Name: file.Name,
		Url:  fmt.Sprintf("%v/%v/%v/%v", os.Getenv("PUBLIC_URL"), controller.STATIC_PATH, file.ID, file.Name),
		Type: file.Type,
		Size: helper.ByteCountSI(file.Size),
	}
}
