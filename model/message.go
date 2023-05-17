package model

import "time"

type Message struct {
	FromId      int         `json:"fromId"`
	TargetId    int         `json:"targetId"`
	SendType    int         `json:"sendType"`
	MessageType int         `json:"messageType"`
	Content     interface{} `json:"content"`
	Time        time.Time   `json:"time"`
}
