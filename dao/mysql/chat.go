package mysql

import (
	"chatchat/model"
	"database/sql"
	"fmt"
)

func InsertIntoMysql(db *sql.DB, message model.Message, key string) error {
	//通过sendtype来判断是私发还是群发
	id := fmt.Sprintf("%sto%s", message.FromId, message.TargetId)
	_, err := db.Exec("insert  into  `messgae` (id,content,time,type,keyss) values (?,?,?,?,?) ", id, message.Content, message.Time, message.SendType, key)

	if err != nil {
		return err
	}
	return nil
}
