package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
    "os"
    "flag"
	"net/http"
    dl "github.com/frzifus/deepL-tg/deepl"
)

func main() {
    pathPtr := flag.String("path", "./data/", "path to config and cert")
    flag.Parse()
    path := *pathPtr
    e, err := os.OpenFile(path + "/e.log", os.O_RDWR | os.O_CREATE | os.O_APPEND,
        0666)
    if err != nil {
        panic(err)
    }
    defer e.Close()

    log.SetOutput(e)
    log.Printf("Path = %s\n", *pathPtr)
    config, err := dl.LoadConfig(path)
    if err != nil {
        panic(err)
    }
	bot, err := tgbotapi.NewBotAPI(config.Token)
    bot.Debug = false

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	w , _:= bot.GetWebhookInfo()
	log.Printf("%s", w.URL)

    log.Printf("Public path: %s", config.Public)
    log.Printf("Private path: %s",  config.Private)

	_, err = bot.SetWebhook(
        tgbotapi.NewWebhookWithCert(config.IP + config.Port + "/" + bot.Token,
            config.Public))
        
	if err != nil {
		log.Fatal(err)
	}

    d := dl.NewDeepL()
    d.SetTargetLang("EN")

	updates := bot.ListenForWebhook("/" + bot.Token)

	go http.ListenAndServeTLS("0.0.0.0:"+config.Port, config.Public,
        config.Private, nil)
	
	for update := range updates {

        d.AddJob(update.Message.Text)
        resp, err := d.Request()
        if err != nil {
            log.Println("fehler beim uebersetzen!")
            d.ResetJobs()
            continue
        }

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp.Translation())
		msg.ReplyToMessageID = update.Message.MessageID
		_, err = bot.Send(msg)
        if err != nil {
            log.Println(err)
        }
        d.ResetJobs()
	}
}
