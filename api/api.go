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
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
)

type Number struct {
	Num1 string `json:"num1"`
	Num2 string `json:"num2"`
	Num3 string `json:"num3"`
	Num4 string `json:"num4"`
	Num5 string `json:"num5"`
	Num6 string `json:"num6"`
	NumB string `json:"numb"`
}

// Index main page
func Index(c echo.Context) error {
	return c.String(http.StatusOK, "hello lotto")
}

// Random return random num 1~45
func Random(c echo.Context) error {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	nums := make(map[int]bool)
	var rn int

	for len(nums) < 6 {
		rn = random.Intn(45) + 1
		nums[rn] = true
	}
	var ret []int
	for i := 1; i <= 45; i++ {
		_, exist := nums[i]
		if exist {
			ret = append(ret, i)
		}
	}

	buf, err := json.Marshal(ret)
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

func GetLottoNum(c echo.Context) error {
	// round := c.Param("round")

	res, err := http.Get("https://dhlottery.co.kr/gameResult.do?method=byWin")
	if err != nil {
		log.Println(err)
		return c.String(http.StatusMethodNotAllowed, "http.Get fail")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("http.Status not OK")
		return c.String(http.StatusMethodNotAllowed, "http Get fail(no 200)")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	win := doc.Find("div.num.win").Contents().FilterFunction(func(i int, s *goquery.Selection) bool {
		return !s.Is("strong")
	}).Text()
	bonus := doc.Find("div.num.bonus").Contents().FilterFunction(func(i int, s *goquery.Selection) bool {
		return !s.Is("strong")
	}).Text()

	winlist := strings.Split(strings.TrimSpace(win), "\n")

	nums := Number{
		strings.TrimSpace(winlist[0]),
		strings.TrimSpace(winlist[1]),
		strings.TrimSpace(winlist[2]),
		strings.TrimSpace(winlist[3]),
		strings.TrimSpace(winlist[4]),
		strings.TrimSpace(winlist[5]),
		strings.TrimSpace(bonus),
	}

	numsJson, err := json.Marshal(nums)
	if err != nil {
		return c.String(http.StatusMethodNotAllowed, "json marshal fail")
	}
	return c.JSONBlob(http.StatusOK, numsJson)
}
