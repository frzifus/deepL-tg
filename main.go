package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	dl "github.com/frzifus/deepL-tg/deepl"
	"gopkg.in/telegram-bot-api.v4"
)

var (
	cfg   *config
	debug *bool
)

func init() {
	cfgFile := flag.String("c", "config.json", "path to config file")
	debug = flag.Bool("d", false, "set to get debug infos")
	flag.Parse()
	log.Printf("Load config from %s\n", *cfgFile)
	var err error
	if cfg, err = loadConfig(*cfgFile); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func main() {
	log.Println("Create a new Telegram bot")
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Println("Couldn't create bot.")
		os.Exit(1)
	}
	bot.Debug = *debug

	log.Println("Authorized on account", bot.Self.UserName)
	w, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Webhook URL:", w.URL)

	_, err = bot.SetWebhook(
		tgbotapi.NewWebhookWithCert(cfg.getWebHookURL(), cfg.Public))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	d := dl.NewDeepL()
	d.SetTargetLang("EN")
	updates := bot.ListenForWebhook("/" + cfg.Token)

	go http.ListenAndServeTLS("0.0.0.0:"+cfg.Port, cfg.Public, cfg.Private, nil)

	for update := range updates {
		if update.Message.Text == "" {
			continue
		}
		d.AddJob(update.Message.Text)
		resp, err := d.Request()
		if err != nil {
			log.Println("error in translation")
			d.ResetJobs()
			continue
		}
		t, err := resp.Translation()
		if err != nil {
			log.Println(err)
			d.ResetJobs()
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, t)
		msg.ReplyToMessageID = update.Message.MessageID
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		d.ResetJobs()
	}
}
