package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/beloslav13/telegram-vk-posts/config"
	"github.com/beloslav13/telegram-vk-posts/logger"
	"net/http"
	"os"
	"time"
)

type Webhook struct {
	StatusWebhook struct {
		Ok     bool `json:"ok"`
		Result struct {
			Url string `json:"url"`
		} `json:"result"`
	}
	SetWebhook struct {
		Ok          bool   `json:"ok"`
		Result      bool   `json:"result"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}
}

// CheckWebhook проверяет веб-хук и если не установлен - устанавливает
func (w *Webhook) CheckWebhook(tgConfig config.Config) (bool, error) {
	url := config.TelegramApiUrl + "/bot" + tgConfig.Telegram.Token + "/getWebhookInfo"
	r, err := http.Get(url)
	if isErr := logger.ForError(err); isErr {
		return false, err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&w.StatusWebhook)
	if isErr := logger.ForError(err); isErr {
		return false, err
	}

	if w.StatusWebhook.Result.Url != "" {
		now := time.Now()
		fmt.Fprintf(os.Stdout, "%s: webhook is already set\n", now.Format("2006-01-02 15:04:05"))
		logger.LogFile.Printf("%s: webhook is already set %+v\n",
			now.Format("2006-01-02 15:04:05"), w.StatusWebhook)
		return true, nil
	} else {
		isSet, err := w.SetWebhooks(tgConfig)
		if isErr := logger.ForError(err); isErr {
			return false, err
		}
		if isSet {
			now := time.Now()
			w.StatusWebhook.Ok = true
			w.StatusWebhook.Result.Url = config.Domain
			fmt.Fprintf(os.Stdout, "%s: webhook was set.\n", now.Format("2006-01-02 15:04:05"))
			logger.LogFile.Printf("%s: webhook was set %+v, url: %s\n",
				now.Format("2006-01-02 15:04:05"), w.SetWebhook, w.StatusWebhook.Result.Url)
			return true, nil
		} else {
			now := time.Now()
			fmt.Fprintf(os.Stdout, "%s: webhook not set.\n", now.Format("2006-01-02 15:04:05"))
			return false, nil
		}
	}

}

// SetWebhooks устанавливает веб-хук
func (w *Webhook) SetWebhooks(tgConfig config.Config) (bool, error) {
	url := config.TelegramApiUrl + "/bot" + tgConfig.Telegram.Token + "/setWebhook?url=" + config.Domain + "/telegram/dialog/"
	r, err := http.Get(url)
	if isErr := logger.ForError(err); isErr {
		return false, err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&w.SetWebhook)
	if isErr := logger.ForError(err); isErr {
		return false, err
	}

	if w.SetWebhook.Result {
		return true, nil
	}
	return false, errors.New(w.SetWebhook.Description)

}
