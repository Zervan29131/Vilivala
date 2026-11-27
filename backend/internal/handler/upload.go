package handler

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadDir is where uploaded images are stored (served at /uploads).
const UploadDir = "uploads"

const maxUploadBytes = 5 << 20

var imageExt = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

// AdminUpload POST /api/admin/upload (multipart field "file")
func AdminUpload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadBytes+(1<<20))
	fh, err := c.FormFile("file")
	if err != nil {
		var mbe *http.MaxBytesError
		if errors.As(err, &mbe) {
			Fail(c, http.StatusRequestEntityTooLarge, 413, "文件不能超过 5MB")
			return
		}
		Fail(c, http.StatusBadRequest, 400, "请通过 file 字段上传文件")
		return
	}
	if fh.Size > maxUploadBytes {
		Fail(c, http.StatusRequestEntityTooLarge, 413, "文件不能超过 5MB")
		return
	}
	src, err := fh.Open()
	if err != nil {
		Fail(c, http.StatusBadRequest, 400, "无法读取文件")
		return
	}
	defer src.Close()

	// detect type from content, never from the file name
	head := make([]byte, 512)
	n, _ := io.ReadFull(src, head)
	ct := http.DetectContentType(head[:n])
	ext, ok := imageExt[ct]
	if !ok {
		Fail(c, http.StatusBadRequest, 400, "仅支持 jpg、png、gif、webp 图片")
		return
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		Fail(c, http.StatusInternalServerError, 500, "保存失败")
		return
	}

	sub := time.Now().Format("200601")
	dir := filepath.Join(UploadDir, sub)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		Fail(c, http.StatusInternalServerError, 500, "保存失败")
		return
	}
	rnd := make([]byte, 8)
	_, _ = rand.Read(rnd)
	name := hex.EncodeToString(rnd) + ext
	dst, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		Fail(c, http.StatusInternalServerError, 500, "保存失败")
		return
	}
	defer dst.Close()
	size, err := io.Copy(dst, src)
	if err != nil {
		Fail(c, http.StatusInternalServerError, 500, "保存失败")
		return
	}
	OK(c, gin.H{"url": "/uploads/" + sub + "/" + name, "size": size, "content_type": ct})
}
