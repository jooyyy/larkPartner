package watcher

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (s *Service) startWatcher() {
	log.Println("start watcher")
	var retryTimes int
	for {
		log.Println("check alive")
		if retryTimes > 3 {
			s.notify()
			retryTimes = 0
			time.Sleep(time.Hour)
			continue
		}

		time.Sleep(s.duration)

		_, err := http.Get(s.link)
		if err != nil {
			log.Println("error", err)
			retryTimes++
			continue
		}

		retryTimes = 0
	}
}

func (s *Service) notify() {
	log.Println("notify")
	var msg = struct {
		MsgType string `json:"msg_type"`
		Content struct {
			Text string `json:"text"`
		} `json:"content"`
	}{
		MsgType: "text",
	}
	msg.Content.Text = "crash"
	raw, _ := json.Marshal(msg)
	if req, err := http.NewRequest(
		"POST",
		s.larkUrl,
		bytes.NewBuffer(raw),
	); err != nil {
		log.Println(err)
		return
	} else {
		client := &http.Client{}
		if resp, err := client.Do(req); err != nil {
			log.Println(err)
		} else {
			defer resp.Body.Close()
		}
	}
}
