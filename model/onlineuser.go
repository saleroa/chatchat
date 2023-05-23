package model

import (
	"github.com/gorilla/websocket"
	"log"
)

type OnLineUser struct {
	UserId       int
	Coon         *websocket.Conn
	ReadChannel  chan Message
	WriteChannel chan Message
}

// read just receive ,and it will send to write
func (ou OnLineUser) Read() {

	for {
		select {
		//要把消息写进这里
		case message, ok := <-ou.ReadChannel:
			if !ok {
				log.Printf("%s ReadService close", ou.UserId)
				return
			}
			ou.WriteChannel <- message
		}
	}
}

// websocket 把消息写进管道内

func (ou OnLineUser) Write() {
	for {
		select {
		case message, ok := <-ou.WriteChannel:
			//close(ou.WriteChannel)
			if !ok {
				log.Printf("%s WriteService close", ou.UserId)
				return
			}
			err := ou.Coon.WriteJSON(message)

			if err != nil {
				log.Printf("%sWrite Err:%s\n", ou.UserId, err.Error())
				return
			}

		}
	}
}
