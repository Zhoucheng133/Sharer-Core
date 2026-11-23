package utils

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/image/draw"

	"github.com/gin-gonic/gin"
)

const previewWidth = 150

// GetPreview 返回宽度不超过 150px 的图片缩略图。
// 非图片或无法解码的图片直接返回源文件，行为与 /api/raw 一致。
func GetPreview(c *gin.Context, basePath string, username string, password string, secret string) {
	decodedPath, ok := resolveFilePath(c, basePath, username, password, secret)
	if !ok {
		return
	}

	file, err := os.Open(decodedPath)
	if err != nil {
		c.JSON(200, gin.H{
			"ok":  false,
			"msg": "File not found",
		})
		return
	}
	defer file.Close()

	// 先做 MIME 探测，避免把非图片内容交给图像解码器。
	header := make([]byte, 512)
	n, err := file.Read(header)
	if err != nil && err != io.EOF {
		c.File(decodedPath)
		return
	}
	if !strings.HasPrefix(http.DetectContentType(header[:n]), "image/") {
		c.File(decodedPath)
		return
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		c.File(decodedPath)
		return
	}

	img, format, err := image.Decode(file)
	if err != nil {
		c.File(decodedPath)
		return
	}

	bounds := img.Bounds()
	if bounds.Dx() <= 0 || bounds.Dx() <= previewWidth {
		c.File(decodedPath)
		return
	}
	height := bounds.Dy() * previewWidth / bounds.Dx()
	if height < 1 {
		height = 1
	}

	thumbnail := image.NewRGBA(image.Rect(0, 0, previewWidth, height))
	draw.CatmullRom.Scale(thumbnail, thumbnail.Bounds(), img, bounds, draw.Src, nil)

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err = jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: 80})
	case "png":
		err = png.Encode(&buf, thumbnail)
	case "gif":
		err = gif.Encode(&buf, thumbnail, nil)
	default:
		c.File(decodedPath)
		return
	}
	if err != nil {
		c.File(decodedPath)
		return
	}

	c.Data(http.StatusOK, "image/"+format, buf.Bytes())
}
