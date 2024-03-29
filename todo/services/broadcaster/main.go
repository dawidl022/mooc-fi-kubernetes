package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/nats-io/nats.go"

	"github.com/dawidl022/mooc-fi-kubernetes/todo/services/broadcaster/config"
)

func main() {
	done := make(chan struct{})
	go broadcastOnMessage()

	<-done
}

func broadcastOnMessage() {
	conf := mustLoadConfig()
	tmpl := mustCreateTemplate(conf.SendMessageTemplate)
	nc := mustInitNatsClient(conf.NatsUrl)

	_, err := nc.QueueSubscribe("todo", "todo-broadcaster", func(msg *nats.Msg) {
		err := broadcastMessage(conf.SendMessageUrl, tmpl, string(msg.Data))
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}

func mustLoadConfig() *config.Config {
	conf, err := config.GetConf()
	if err != nil {
		log.Fatal(err)
	}
	return conf
}

func mustCreateTemplate(templateString string) *template.Template {
	tmpl, err := template.New("messageTemplate").Parse(templateString)
	if err != nil {
		log.Fatal(err)
	}
	return tmpl
}

func mustInitNatsClient(url string) *nats.Conn {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	return nc
}

func broadcastMessage(url string, tmpl *template.Template, msg string) error {
	escapedMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	tmpl.Execute(&buf, struct{ Message string }{string(escapedMsg)})
	_, err = http.Post(url, "application/json", &buf)
	return err
}
