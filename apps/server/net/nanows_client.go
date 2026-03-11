package net

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/appditto/pippin_nano_wallet/libs/log"
	"github.com/appditto/pippin_nano_wallet/libs/wallet"
	"github.com/recws-org/recws"
)

type wsSubscribe struct {
	Action  string              `json:"action"`
	Topic   string              `json:"topic"`
	Ack     bool                `json:"ack"`
	Id      string              `json:"id"`
	Options map[string][]string `json:"options"`
}

type ConfirmationResponse struct {
	Topic   string                 `json:"topic"`
	Time    string                 `json:"time"`
	Message map[string]interface{} `json:"message"`
}

type WSCallbackBlock struct {
	Type           string `json:"type"`
	Account        string `json:"account"`
	Previous       string `json:"previous"`
	Representative string `json:"representative"`
	Balance        string `json:"balance"`
	Link           string `json:"link"`
	LinkAsAccount  string `json:"link_as_account"`
	Work           string `json:"work"`
	Signature      string `json:"signature"`
	Destination    string `json:"destination"`
	Source         string `json:"source"`
	Subtype        string `json:"subtype"`
}

type WSCallbackMsg struct {
	IsSend  string          `json:"is_send"`
	Block   WSCallbackBlock `json:"block"`
	Account string          `json:"account"`
	Hash    string          `json:"hash"`
	Amount  string          `json:"amount"`
}

func StartNanoWSClient(wsUrl string, callbackChan *chan *WSCallbackMsg, w *wallet.NanoWallet, newAccountChan <-chan string) {
	log.Infof("Starting StartNanoWSClient")
	ctx, cancel := context.WithCancel(context.Background())
	sentSubscribe := false
	ws := recws.RecConn{}

	addresses, err := w.GetAllAccountAddresses()
	if err != nil {
		addresses = []string{}
	}

	log.Infof("Subscribed to %d accounts", len(addresses))

	subRequest := wsSubscribe{
		Action: "subscribe",
		Topic:  "confirmation",
		Ack:    false,
		Options: map[string][]string{
			"accounts": addresses,
		},
	}
	ws.Dial(wsUrl, nil)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer func() {
		signal.Stop(sigc)
		cancel()
	}()

	// Goroutine to handle new account addresses
	go func() {
		for newAccount := range newAccountChan {
			log.Infof("Resubscribing with account: %s", newAccount)
			//todo: this should use action: "update", "accounts_add" from v21 but node rpc proxies don't support :(
			addresses = append(addresses, newAccount)
			// Send update message to WebSocket
			updateRequest := wsSubscribe{
				Action: "subscribe",
				Topic:  "confirmation",
				Ack:    false,
				Options: map[string][]string{
					"accounts": addresses,
				},
			}
			if err := ws.WriteJSON(updateRequest); err != nil {
				log.Infof("Error sending update request %s", ws.GetURL())
			} else {
				log.Infof("Successfully sent update request for new account %s", newAccount)
			}
		}
	}()

	for {
		select {
		case <-sigc:
			cancel()
			return
		case <-ctx.Done():
			go ws.Close()
			log.Infof("Websocket closed %s", ws.GetURL())
			return
		default:
			if !ws.IsConnected() {
				sentSubscribe = false
				log.Infof("Websocket disconnected %s", ws.GetURL())
				time.Sleep(2 * time.Second)
				continue
			}

			// Sent subscribe with ack
			if !sentSubscribe {
				if err := ws.WriteJSON(subRequest); err != nil {
					log.Infof("Error sending subscribe request %s", ws.GetURL())
					time.Sleep(2 * time.Second)
					continue
				} else {
					sentSubscribe = true
				}
			}

			var confMessage ConfirmationResponse
			err := ws.ReadJSON(&confMessage)
			if err != nil {
				log.Infof("Error: ReadJSON %s", ws.GetURL())
				sentSubscribe = false
				continue
			}

			// Trigger callback
			if confMessage.Topic == "confirmation" {
				log.Infof("Received websocket confirmation")
				var deserialized WSCallbackMsg
				serialized, err := json.Marshal(confMessage.Message)
				if err != nil {
					log.Infof("Error: Marshal ws %v", err)
					continue
				}
				if err := json.Unmarshal(serialized, &deserialized); err != nil {
					log.Errorf("Error: decoding the callback to WSCallbackMsg %v", err)
					continue
				}
				*callbackChan <- &deserialized
			}
		}
	}
}
