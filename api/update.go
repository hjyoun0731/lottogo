package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

// UploadFile upload apk to server
func UploadFile(c echo.Context) error {
	// fileSize : file size by formvalue
	fileSize, err := strconv.Atoi(c.FormValue("size"))
	if err != nil {
		return err
	}

	// file : file data
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusNotFound, "UploadFile Fail.")
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// file size check like cksum
	var buf bytes.Buffer

	bufSize, err := buf.ReadFrom(src)
	if err != nil {
		return err
	}
	_, err = src.Seek(0, 0)
	if err != nil {
		return err
	}
	if int(bufSize) != fileSize {
		return c.String(http.StatusMethodNotAllowed, "file size ("+strconv.Itoa(int(bufSize))+") is not matched.")
	}

	// file save
	dst, err := os.Create("./files/" + file.Filename)
	if err != nil {
		fmt.Println(err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "file upload success.")
}

func DownloadFile(c echo.Context) error {
	var apks string = "lottogo.apks"

	_, err := os.Stat(apks)
	if err != nil {
		return c.String(http.StatusMethodNotAllowed, "can't check fileinfo")
	} else if os.IsNotExist(err) {
		return c.String(http.StatusMethodNotAllowed, "file not exist")
	}
	return c.Attachment("./files/"+apks, apks)
}
