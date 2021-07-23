package server

import (
	"encoding/json"
	"fmt"
	"github.com/beloslav13/telegram-vk-posts/pkg/telegram/models"
	"github.com/beloslav13/telegram-vk-posts/redis"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	var respTelegram models.RespTelegram
	err := json.NewDecoder(r.Body).Decode(&respTelegram)
	if err != nil {
		fmt.Println(err)
		return
	}
	//w.Write([]byte(`{"status": 200`))
	fmt.Printf("%+v\n", respTelegram)
	val, err := models.Rdb.Client.Get(redis.Ctx, "state").Result()
	if err != nil {
		fmt.Println(models.Response(&respTelegram, respTelegram.Message.Text))
		return
	}

	models.Response(&respTelegram, val)

	fmt.Println("state", val)

}

func StartServer() (bool, error) {
	http.HandleFunc("/telegram/dialog/", index) // each request calls handler
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		return false, err
	}
	return true, nil
}
