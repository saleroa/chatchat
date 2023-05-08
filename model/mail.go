package model

type Mail struct {
	Username string `form:"username" json:"username" `
	MailID   string `form:"mailID" json:"mailID"`
}
