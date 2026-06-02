package utilit

import (
	"log"
	"net/http"
	"net/url"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init() *tgbotapi.BotAPI {
	token := os.Getenv("BOT_TOKEN")
	proxy := os.Getenv("PROXY")
	if token == "" {
		log.Panic("BOT_TOKEN не задан!")
	}

	// Если есть авторизация: "http://логин:пароль@ip_адрес:порт"
	proxyURLStr := proxy

	proxyURL, err := url.Parse(proxyURLStr)
	if err != nil {
		log.Panicf("Неверный формат URL прокси: %v", err)
	}

	// Настраиваем транспорт с указанием прокси
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	bot, err := tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, httpClient)
	if err != nil {
		log.Panic(err)
	}

	return bot
}
