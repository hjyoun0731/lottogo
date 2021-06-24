package api

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
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

func CreateTable(c echo.Context) error {
	err := NewTable()
	if err != nil {
		return c.String(http.StatusMethodNotAllowed, "table create fail")
	}
	return c.String(http.StatusOK, "table created")
}

func GetLottoNum(c echo.Context) error {
	lottoUrl := "https://dhlottery.co.kr/gameResult.do?method=byWin"

	round := c.Param("round")
	if round != "latest" {
		lottoUrl = lottoUrl + "&drwNo=" + round
	}

	res, err := http.Get(lottoUrl)
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

	var nums []string
	for i := 0; i < 6; i++ {
		nums = append(nums, strings.TrimSpace(winlist[i]))
	}
	nums = append(nums, strings.TrimSpace(bonus))

	numsJson, err := json.Marshal(nums)
	if err != nil {
		return c.String(http.StatusMethodNotAllowed, "json marshal fail")
	}
	return c.JSONBlob(http.StatusOK, numsJson)
}
