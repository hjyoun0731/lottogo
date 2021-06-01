package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Number struct {
	Num1 int `json:"num1"`
}

// Index main page
func Index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	_, err := fmt.Fprint(w, "Welcome!\n")
	if err != nil {
		log.Println("router.Index error")
	}
}

// Random return random num 1~45
func Random(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	nums := &Number{}

	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	nums.Num1 = random.Intn(45) + 1

	buf, err := json.Marshal(nums)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(buf)
}

// UploadFile upload apk to server
func UploadFile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseMultipartForm(0 << 100)
	if err != nil {
		log.Println("UploadFile ParseMultipartForm fail")
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(404)

		log.Println("UploadFile Fail.")
		return
	}
	defer file.Close()

	var buf bytes.Buffer

	bufSize, err := buf.ReadFrom(file)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	formSize, err := strconv.Atoi(r.FormValue("size"))
	if err != nil {
		log.Println(err)
	}
	if int(bufSize) != formSize {
		w.WriteHeader(405)
		log.Println("Upload file fail - size")
		return
	}

	upFile, err := os.Create("./files/" + time.Now().String() + "____" + handler.Filename)
	if err != nil {
		fmt.Println(err)
	}
	defer upFile.Close()

	_, err = upFile.Write(buf.Bytes())
	if err != nil {
		log.Println("UploadFile fileBytes write fail")
	}
	log.Println("UploadFile success")
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db := NewDb()
	defer CloseDb(db)
	rows, err := db.Query("SELECT user_id, user_name, user_password, created, updated  FROM User_Table where user_id >= ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var ui UserInfo
	for rows.Next() {
		err := rows.Scan(&ui.UserId, &ui.UserName, &ui.UserPassword, &ui.Created, &ui.Updated)
		if err != nil {
			log.Fatal(err)
		}
	}

	buff, err := json.Marshal(ui)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(buff)
}

func InsertUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	ui := &UserInfo{
		UserName: r.FormValue("user_name"),
		UserPassword: r.FormValue("user_password"),
		Created: r.FormValue("created"),
		Updated: r.FormValue("updated"),
	}

	db := NewDb()
	defer CloseDb(db)
	_, err := db.Exec("INSERT INTO User_Table(user_name, user_password, created, updated) VALUES(?, ?, ?, ?)", ui.UserName, ui.UserPassword, ui.Created, ui.Updated)
	if err != nil {
		log.Fatal(err)
	}
	w.Write([]byte("insert successed"))

}
