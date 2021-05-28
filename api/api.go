package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	upFile, err := os.Create("./files/" + handler.Filename)
	if err != nil {
		fmt.Println(err)
	}
	defer upFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	_, err = upFile.Write(fileBytes)
	if err != nil {
		log.Println("UploadFile fileBytes write fail")
	}

	filesize, err := os.Stat("./files/" + handler.Filename)
	if err != nil {
		log.Println(err)
	}

	log.Println("size:", filesize.Size())
	formSize, err := strconv.Atoi(r.FormValue("size"))
	if err != nil {
		log.Println(err)
	}
	if formSize != int(filesize.Size()) {
		w.WriteHeader(405)
	}

	log.Println("UploadFile success")
}
