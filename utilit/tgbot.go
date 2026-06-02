package utilit

import (
	"encoding/base64" // Добавлено для авторизации прокси
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

	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Panicf("Неверный формат URL прокси: %v", err)
	}

	// 1. Создаем транспорт и указываем прокси-сервер
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	// 2. Если в URL прокси переданы логин и пароль, настраиваем авторизацию
	if proxyURL.User != nil {
		password, _ := proxyURL.User.Password()
		username := proxyURL.User.Username()

		// Кодируем данные в формат Base64 для заголовка Basic Auth
		auth := username + ":" + password
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

		// Передаем заголовок авторизации в прокси-туннель (для HTTPS/CONNECT запросов)
		transport.ProxyConnectHeader = http.Header{
			"Proxy-Authorization": []string{basicAuth},
		}
	}

	// 3. Собираем итоговый HTTP-клиент
	httpClient := &http.Client{
		Transport: transport,
	}

	bot, err := tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, httpClient)
	if err != nil {
		log.Panic(err)
	}

	return bot
}
