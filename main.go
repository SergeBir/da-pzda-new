package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var answers = map[string]string{
	"да":       "пизда",
	"пизда":    "да",
	"нет":      "пидора ответ",
	"здрасьте": "забор покрасьте",
	"300":      "отсоси у тракториста",
	"триста":   "отсоси у тракториста",
	"точно":    "соси сочно",
}

func main() {
	token, err := ioutil.ReadFile("token.txt")
	if err != nil {
		log.Printf("File reading error: %v", err)
		return
	}

	bot, err := tgbotapi.NewBotAPI(string(token))
	if err != nil {
		log.Fatalf("Error connecting to the bot: %v", err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	upd, err := bot.GetUpdatesChan(ucfg)
	if err != nil {
		log.Fatalf("Error getting updates channel: %v", err)
	}
	time.Sleep(time.Millisecond * 500)
	upd.Clear()

	for {
		select {
		case update := <-upd:
			if update.Message != nil { // make sure the incoming update is a message
				if reply, ok := answers[strings.TrimRight(strings.ToLower(update.Message.Text), "аеёиоуыэюя")]; ok {
					echoMsg := tgbotapi.NewMessage(update.Message.Chat.ID, addTailVowels(reply, update.Message.Text))
					echoMsg.ReplyToMessageID = update.Matcher.MessageID
					log.Printf("Sending %s", reply)
					_, err := bot.Send(echoMsg)
					if err != nil {
						log.Printf("Error sending message: %v", err)
					}
				}
			}
		}
	}
}

func addTailVowels(base string, original string) string {
	tail := findVowelTail(original)
	return base + tail
}

func findVowelTail(s string) string {
	r, _ := regexp.Compile("[аеёиоуыэюя]+$")
	found := r.FindString(s)
	return found
}
