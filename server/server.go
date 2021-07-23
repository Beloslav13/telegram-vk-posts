package server

import (
	"encoding/json"
	"fmt"
	"github.com/beloslav13/telegram-vk-posts/pkg/telegram/models"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	var respTelegram models.RespTelegram
	err := json.NewDecoder(r.Body).Decode(&respTelegram)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", respTelegram)
	//w.Write([]byte(`{"status": 200`))
	models.Response(&respTelegram, respTelegram.Message.Text)

}

func StartServer() (bool, error) {
	http.HandleFunc("/telegram/dialog/", index) // each request calls handler
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		return false, err
	}
	return true, nil
}
