package models

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

var States = map[string]func() string{
	"/help":  Help,
	"/posts": GetPosts,
}

func Response(command string) string {
	state, ok := States[command]
	if ok {
		return state()
	} else {
		return notFoundState()
	}
}

func Help() string {
	return "help!"
}

func GetPosts() string {
	return "GetPosts!"
}

func notFoundState() string {
	return "No found state."
}
