package mysql

import (
	"chatchat/model"
	"database/sql"
)

func InsertIntoMysql(db *sql.DB, message model.Message, key string) error {
	_, err := db.Exec("insert  into  `message` (mid,fromid,toid,content,time,sendtype,messagetype) values (?,?,?,?,?,?,?) ", 0, message.FromId, message.TargetId, message.Content, message.Time, message.SendType, message.MessageType)

	if err != nil {
		return err
	}
	return nil
}
