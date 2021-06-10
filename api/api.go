package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type Number struct {
	Num1 int `json:"num1"`
	Num2 int `json:"num2"`
	Num3 int `json:"num3"`
	Num4 int `json:"num4"`
	Num5 int `json:"num5"`
	Num6 int `json:"num6"`
}

// Index main page
func Index(c echo.Context) error {
	return c.String(http.StatusOK, "hello lotto")
}

// Random return random num 1~45
func Random(c echo.Context) error {
	nums := &Number{}

	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	nums.Num1 = random.Intn(45) + 1
	nums.Num2 = random.Intn(45) + 1
	nums.Num3 = random.Intn(45) + 1
	nums.Num4 = random.Intn(45) + 1
	nums.Num5 = random.Intn(45) + 1
	nums.Num6 = random.Intn(45) + 1

	buf, err := json.Marshal(nums)
	if err != nil {
		log.Println(err)
	}

	return c.JSONBlob(http.StatusOK, buf)
}

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

func GetUserInfo(c echo.Context) error {
	db := NewDb()
	defer CloseDb(db)
	rows, err := db.Query("SELECT id, name, password, created, updated  FROM user_table where id >= ?", 1)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ui UserInfo
	for rows.Next() {
		err := rows.Scan(&ui.Id, &ui.Name, &ui.Password, &ui.Created, &ui.Updated)
		if err != nil {
			return err
		}
	}

	buff, err := json.Marshal(ui)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(buff))
}

func NewUserInfo(c echo.Context) error {

	ui := &UserInfo{
		Name:     c.FormValue("name"),
		Password: c.FormValue("password"),
	}

	db := NewDb()
	defer CloseDb(db)
	_, err := db.Exec("INSERT INTO user_table(name, password) VALUES(?, ?)", ui.Name, ui.Password)
	if err != nil {
		log.Println(err)
	}
	return c.String(http.StatusOK, "SignUp Success.")
}

func SignIn(c echo.Context) error {
	db := NewDb()
	defer CloseDb(db)
	params := make(map[string]string)
	var id int

	name := c.FormValue("name")
	password := c.FormValue("password")
	id, dbPassword := queryPw(db, name)

	if password != dbPassword {
		params["pwd"] = "no match"
		_ = c.Bind(&params)

		return c.JSON(http.StatusMethodNotAllowed, params["pwd"])
	}

	accessToken, err := generateToken(c, id, name)
	if err != nil {
		params["token"] = fmt.Sprint(err)
		return c.JSON(http.StatusMethodNotAllowed, params["token"])
	}

	c.Response().Header().Set("Cache-Control", "no-store no-cache, must-revalidate, post-check=0, pre-check=0, max-age=0")
	c.Response().Header().Add("Last-Modified", time.Now().String())
	c.Response().Header().Add("pragma", "no-cache")
	// c.Response().Header().Add("Expires", "-1")
	cookie := new(http.Cookie)
	cookie.Name = "access-token"
	cookie.Value = accessToken
	cookie.Expires = time.Now().Add(ExpirationTime)
	c.SetCookie(cookie)
	params["token"] = accessToken

	return c.JSON(http.StatusOK, params["token"])
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

func CreateTable(c echo.Context) error {
	err := NewTable()
	if err != nil {
		return c.String(http.StatusMethodNotAllowed, "table create fail")
	}
	return c.String(http.StatusOK, "table created")
}
