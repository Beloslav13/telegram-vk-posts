package models

import (
	"errors"
	"fmt"
	"github.com/beloslav13/telegram-vk-posts/redis"
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
	token, exists := os.LookupEnv("TOKEN")

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
	"/start":       Start,
	"/help":        Help,
	"/posts":       Posts,
	"getPosts":     GetPosts,
	"getPostsWait": GetPostsWait,
	"/exit":        Exit,
}

var Rdb, _ = redis.NewDatabase("localhost:6379")

func Response(r *RespTelegram, command string) string {
	state, ok := States[command]
	if ok {
		return state(r)
	} else {
		return notFoundState(r)
	}
}

func Start(r *RespTelegram) string {
	tmp := "Привет! В этом боте можно получить публикации из Вконтакте.\nВоспользуйся командой /posts и просто следуй инструкции."
	err := r.sendMessage(strconv.Itoa(int(r.Message.Chat.Id)), tmp)
	if err != nil {
		fmt.Println(err)
		return "Err"
	}
	fmt.Println("Message text>>>>", r.Message.Text)
	//Rdb.Client.Set(redis.Ctx, "state", r.Message.Text, 0)
	return "Command is start...."
}

func Help(r *RespTelegram) string {
	tmp := "В этом боте можно получить публикации из Вконтакте, воспользуйтесь /posts"
	err := r.sendMessage(strconv.Itoa(int(r.Message.Chat.Id)), tmp)
	if err != nil {
		fmt.Println(err)
		return "Err"
	}
	Rdb.Client.Del(redis.Ctx, "state")
	return "help!"
}

func Posts(r *RespTelegram) string {
	if r.Message.Text == "/exit" {
		Exit(r)
		return "Вышли..."
	}

	Rdb.Client.Set(redis.Ctx, "state", "getPosts", 0)
	tmp := "Пришлите ссылку группы Вконтакте."
	err := r.sendMessage(strconv.Itoa(int(r.Message.Chat.Id)), tmp)
	if err != nil {
		fmt.Println(err)
		return "Err"
	}
	return "Posts!"
}

func GetPosts(r *RespTelegram) string {
	if r.Message.Text == "/exit" {
		Exit(r)
		val, _ := Rdb.Client.Get(redis.Ctx, "state").Result()
		fmt.Println(val)
		return "Вышли..."
	}

	tmp := "Отсылаю публикации Вконакте....\n Спарсить ещё что-нибудь?"
	err := r.sendMessage(strconv.Itoa(int(r.Message.Chat.Id)), tmp)
	if err != nil {
		fmt.Println(err)
		return "Err"
	}
	//Rdb.Client.Del(redis.Ctx, "state")
	Rdb.Client.Set(redis.Ctx, "state", "getPostsWait", 0)
	return ""
}

func GetPostsWait(r *RespTelegram) string {
	if r.Message.Text == "да" {
		GetPosts(r)
		return "Парсим еще раз"
	} else if r.Message.Text == "нет" {
		fmt.Println("Ответ нет...")
		Exit(r)
		return "Выходим...."
		//Rdb.Client.Set(redis.Ctx, "state", "/exit", 0)
	} else { // Поправить
		tmp := "Ничего не понял! Воспользуйся командой /help"
		err := r.sendMessage(strconv.Itoa(int(r.Message.Chat.Id)), tmp)
		if err != nil {
			fmt.Println(err)
			return "Err"
		}
		Rdb.Client.Set(redis.Ctx, "state", "/help", 0)
		return "Не понятно"
	}
}

func Exit(r *RespTelegram) string {
	tmp := "Пока-пока!"
	err := r.sendMessage(strconv.Itoa(int(r.Message.Chat.Id)), tmp)
	if err != nil {
		fmt.Println(err)
		return "Err"
	}
	Rdb.Client.Del(redis.Ctx, "state")
	return "Вышли"
}

func notFoundState(r *RespTelegram) string {
	tmp := "Воспользуйтесь командой /help"
	err := r.sendMessage(strconv.Itoa(int(r.Message.Chat.Id)), tmp)
	if err != nil {
		fmt.Println(err)
		return "Err"
	}
	return "No found state."
}
