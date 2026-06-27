package controllers

import (
	"blog/utils"
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type FileUploadController struct {
	UploadPath string
}

func (f *FileUploadController) UploadFile(ctx *gin.Context) {

	// 获取上传的文件
	file, err := ctx.FormFile("blogfile")
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	// 保存文件
	filename := filepath.Base(file.Filename)
	dst := filepath.Join(f.UploadPath, filename)
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		utils.InternalServerError(ctx, err.Error())
		return
	}

	utils.Success(ctx, gin.H{
		"url": fmt.Sprintf("/uploads/%s", filename),
	})
}
