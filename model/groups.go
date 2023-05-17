package model

import "time"

type Group struct {
	Id       int //后来加
	Name     string
	Time     time.Time
	MangerID int
}
