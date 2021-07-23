package models

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type RespTelegram struct {
	UpdateId int32   `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageId int32  `json:"message_id"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

type User struct {
	Id        int32  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Chat struct {
	Id        int32  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (r *RespTelegram) sendMessage(chatId string, text string) error {
	token, exists := os.LookupEnv("TOKENS")

	if !exists {
		return errors.New("Token not found.")
	}
	tgUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?", token)
	u := url.Values{}
	u.Set("chat_id", chatId)
	u.Set("text", text)

	req := tgUrl + u.Encode()
	fmt.Println(req)
	resp, err := http.Get(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil

}

var States = map[string]func(r *RespTelegram) string{
	"/help":  Help,
	"/posts": GetPosts,
}

func Response(r *RespTelegram, command string) string {
	state, ok := States[command]
	if ok {
		return state(r)
	} else {
		return notFoundState()
	}
}

func Help(r *RespTelegram) string {
	return "help!"
}

func GetPosts(r *RespTelegram) string {
	tmp := "Пришлите ссылку группы Вконтакте."
	err := r.sendMessage(strconv.Itoa(int(r.Message.Chat.Id)), tmp)
	if err != nil {
		fmt.Println(err)
		return "Err"
	}
	return "GetPosts!"
}

func notFoundState() string {
	return "No found state."
}
