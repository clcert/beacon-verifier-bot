package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type response struct {
	Pulse       string            `json:"pulse"`
	Valid       bool              `json:"valid"`
	CheckedDate string            `json:"checked_date"`
	Sources     map[string]result `json:"Sources"`
}

type result struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason"`
}

const interval = 1 * time.Minute

func main() {
	token := os.Getenv("TG_TOKEN")
	debug := os.Getenv("DEBUG") == "1"
	group, err := strconv.ParseInt(os.Getenv("TG_GROUP_ID"), 10, 64)
	if err != nil {
		log.Panic(err)
	}
	verifierURL := os.Getenv("BEACON_VERIFIER_API")
	ignoredSources := strings.Fields(os.Getenv("IGNORED_SOURCES"))
	ignoredMap := make(map[string]struct{})
	for _, source := range ignoredSources {
		ignoredMap[source] = struct{}{}
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Starting bot with api %s in group id=%d...", verifierURL, group)
	bot.Debug = true
	for {
		startLoop := time.Now()
		r, err := http.Get(verifierURL)
		if err != nil {
			log.Printf("%+v", err)
		} else {
			respBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Printf("%+v", err)
			} else {
				var resp response
				json.Unmarshal(respBytes, &resp)
				// query to api latest
				badSources := make(map[string]result, 0)
				for source, result := range resp.Sources {
					if _, isIgnored := ignoredMap[source]; !(result.Valid || isIgnored) {
						badSources[source] = result
					}
				}
				if len(badSources) > 0 {
					msgText := fmt.Sprintf("failed validating apis at pulse %s", resp.Pulse)
					for source, result := range badSources {
						msgText += fmt.Sprintf("\n%s: %s", source, result.Reason)
					}
					msg := tgbotapi.NewMessage(group, msgText)
					msg.DisableNotification = true
					resp, err := bot.Send(msg)
					if err != nil {
						log.Printf("%+v", err)
					} else {
						log.Printf("response: %+v", resp)
					}
				} else {
					log.Printf("Everything OK for pulse %s!", resp.Pulse)
					if debug {
						msg := tgbotapi.NewMessage(group, fmt.Sprintf("DEBUG: All sources are validated"))
						msg.DisableNotification = true
						resp, err := bot.Send(msg)
						if err != nil {
							log.Printf("%+v", err)
						} else {
							log.Printf("response: %+v", resp)
						}
					}

				}
			}
		}
		now := time.Now()
		seconds := time.Duration(math.Round(now.Sub(startLoop).Seconds())) * time.Second
		sleepTime := interval - seconds
		log.Printf("Finished this cycle. Sleeping %s...", sleepTime)
		time.Sleep(sleepTime)
	}
}
